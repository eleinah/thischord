package audio

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/eleinah/thischord/internal/helpers"
	"github.com/eleinah/thischord/internal/state"
	"github.com/lrstanley/go-ytdlp"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	ffmpegBufferSize int     = 16384
	maxClampValue    float64 = 32767
	minClampValue    float64 = -32767
)

func playAudio(interactionState *state.InteractionState, url string, stop chan bool, pauseCh chan bool, done chan bool) {
	defer close(done)

	var vc *discordgo.VoiceConnection
	var err error
	var ok bool

	if !state.VoiceConnected {
		vc, err = helpers.JoinUserVoiceChannel(interactionState)
		if err != nil {
			slog.Error("Failed to join voice channel", "error", err.Error())
			interactionState.Reply("Failed to join voice channel.")
			return
		}
	} else {
		vc, ok = interactionState.Session.VoiceConnections[interactionState.GuildID]
		if !ok {
			slog.Error("Error getting voice connection", "exists", err)
			interactionState.Reply("Error with voice connection.")
			return
		}
	}

	songDone := make(chan bool)

	go func() {
		PlayURL(vc, url, stop, pauseCh)
		close(songDone)
	}()

	<-songDone
	slog.Info("Song finished playing")
}

func PlayURL(v *discordgo.VoiceConnection, url string, stop <-chan bool, pauseCh <-chan bool) {
	// TODO: URL validation

	ytOut, err := ytdlp.New().
		NoPlaylist().
		ExtractAudio().
		AudioQuality("0").
		Output("-").
		Run(context.Background(), url)
	if err != nil {
		slog.Error("Failed to create youtube downloader", "url", url, "error", err.Error())
		return
	}

	ytData := bytes.NewReader([]byte(ytOut.Stdout))
	ytOut.Stdout = ""

	pr, pw := io.Pipe()

	go func() {
		err = ffmpeg.Input("pipe:0").
			WithInput(ytData).
			Output("pipe:1",
				ffmpeg.KwArgs{
					"f":      "s16le",
					"ar":     strconv.Itoa(FrameRate),
					"ac":     strconv.Itoa(Channels),
					"acodec": "pcm_s16le",
				}).
			WithOutput(pw).
			Run()
		if err != nil {
			slog.Error("Failed to transcode audio to PCM", "url", url, "error", err.Error())
			pwErr := pw.CloseWithError(err)
			if pwErr != nil {
				slog.Error("Failed to close pipe", "url", url, "error", pwErr.Error())
				return
			}
			return
		}
		pwErr := pw.Close()
		if pwErr != nil {
			slog.Error("Failed to close pipe", "url", url, "error", pwErr.Error())
			return
		}
	}()

	ffmpegbuf := bufio.NewReaderSize(pr, ffmpegBufferSize)

	time.Sleep(100 * time.Millisecond)

	err = v.Speaking(true)
	if err != nil {
		slog.Error("Failed to set speaking to true", "url", url, "error", err.Error())
		return
	}
	defer func() {
		err = v.Speaking(false)
		if err != nil {
			slog.Error("Failed to set speaking to false", "url", url, "error", err.Error())
		}
	}()

	send := make(chan []int16, 2)
	defer close(send)

	closeCh := make(chan bool)
	go func() {
		SendPCM(v, send)
		closeCh <- true
	}()

	minPlayTime := time.NewTimer(500 * time.Millisecond)
	defer minPlayTime.Stop()

	dataRecv := false

	isPaused := false

	for {
		select {
		case newState := <-pauseCh:
			isPaused = newState
			continue
		default:
		}

		if isPaused {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		audiobuf := make([]int16, FrameSize*Channels)

		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			if !dataRecv {
				select {
				case <-minPlayTime.C:
					return
				case <-closeCh:
					return
				}
			}
			return
		}
		if err != nil {
			slog.Error("Failed to read audio data", "url", url, "error", err.Error())
			return
		}

		dataRecv = true

		state.VolumeMtx.Lock()
		currentVol, ok := state.Volume[v.GuildID]
		if !ok {
			currentVol = 1.0
			state.Volume[v.GuildID] = 1.0
		}
		state.VolumeMtx.Unlock()

		for i := range audiobuf {
			newVal := float64(audiobuf[i]) * currentVol

			if newVal > maxClampValue {
				newVal = maxClampValue
			} else if newVal < minClampValue {
				newVal = minClampValue
			}
			audiobuf[i] = int16(newVal)
		}

		select {
		case send <- audiobuf:
		case <-closeCh:
			return
		}
	}
}
