package commands

import "github.com/eleinah/thischord/internal/state"

func Ping(interactionState *state.InteractionState) {
	interactionState.Reply("Pong!")
}
