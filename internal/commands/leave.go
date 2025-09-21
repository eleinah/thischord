package commands

import (
	"github.com/eleinah/thischord/internal/helpers"
	"github.com/eleinah/thischord/internal/state"
)

func Leave(interactionState *state.InteractionState) {
	helpers.LeaveUserVoiceChannel(interactionState)
}
