package commands

import "github.com/eleinah/thischord/internal/state"

func Greet(interactionState *state.InteractionState) {
	interactionState.Reply("Hello, world!")
}
