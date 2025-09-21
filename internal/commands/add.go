package commands

import (
	"fmt"

	"github.com/eleinah/thischord/internal/state"
)

func Add(interactionState *state.InteractionState) {
	opts := interactionState.Interaction.ApplicationCommandData().Options
	num1 := opts[0].IntValue()
	num2 := opts[1].IntValue()
	sum := num1 + num2
	sumString := fmt.Sprintf("%d + %d = %d", num1, num2, sum)

	interactionState.Reply(sumString, false)
}
