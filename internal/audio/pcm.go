package audio

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/eleinah/thischord/internal/state"
	"gopkg.in/hraban/opus.v2"
)

const (
	FrameRate int = 48000
	Channels  int = 2
	FrameSize int = 960
	MaxBytes      = (FrameSize * 2) * 2
)

func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	var err error

	state.OpusEncoder, err = opus.NewEncoder(FrameRate, Channels, opus.AppAudio)
	if err != nil {
		slog.Error("Error creating Opus encoder", "error", err)
		return
	}

	for {
		recv, ok := <-pcm
		if !ok {
			return
		}

		data := make([]byte, MaxBytes)

		if _, err := state.OpusEncoder.Encode(recv, data); err != nil {
			slog.Error("Error encoding PCM data", "error", err)
			return
		}

		if !v.Ready || v.OpusSend == nil {
			return
		}

		v.OpusSend <- data
	}
}
