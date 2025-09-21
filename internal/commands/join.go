package commands

import (
	"github.com/eleinah/thischord/internal/helpers"
	"github.com/eleinah/thischord/internal/state"
)

func Join(interactionState *state.InteractionState) {
	helpers.JoinUserVoiceChannel(interactionState)
}
