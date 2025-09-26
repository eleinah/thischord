package audio

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/eleinah/thischord/internal/state"
)

func ProcessQueue(interactionState *state.InteractionState) {
	go func() {
		for {
			state.QueueMtx.Lock()
			if len(state.Queue[interactionState.GuildID]) == 0 {
				state.PlayingMtx.Lock()
				state.Playing[interactionState.GuildID] = false
				state.PlayingMtx.Unlock()
				state.QueueMtx.Unlock()

				time.Sleep(500 * time.Millisecond)

				vc, ok := interactionState.Session.VoiceConnections[interactionState.GuildID]
				if ok {
					err := vc.Speaking(false)
					if err != nil {
						slog.Error("Error setting speaking to false", "error", err.Error())
						return
					}
					err = vc.Disconnect()
					if err != nil {
						slog.Error("Error disconnecting from voice", "error", err.Error())
						return
					}
				}
				break
			}

			currentURL := state.Queue[interactionState.GuildID][0]
			state.Queue[interactionState.GuildID] = state.Queue[interactionState.GuildID][1:]
			songLength := len(state.Queue[interactionState.GuildID])
			state.QueueMtx.Unlock()

			slog.Info("Playing song", "songs_remaining", strconv.Itoa(songLength))
			interactionState.Reply(fmt.Sprintf("Now playing: %s", currentURL))

			state.StopMtx.Lock()
			stop := make(chan bool)
			state.StopChans[interactionState.GuildID] = stop
			state.StopMtx.Unlock()

			state.PauseChMtx.Lock()
			pauseCh := make(chan bool)
			state.PauseChans[interactionState.GuildID] = pauseCh
			state.PauseChMtx.Unlock()

			state.PausedMtx.Lock()
			pauseCh <- state.Paused[interactionState.GuildID]
			state.PausedMtx.Unlock()

			done := make(chan bool)
			go playAudio(interactionState, currentURL, stop, pauseCh, done)
			<-done

			slog.Info("Song finished, moving to next in queue (if possible)")

			state.PauseChMtx.Lock()
			delete(state.PauseChans, interactionState.GuildID)
			state.PauseChMtx.Unlock()
		}
	}()
}
