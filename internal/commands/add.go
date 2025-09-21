package commands

import (
	"fmt"
	"strconv"

	"github.com/eleinah/thischord/internal/state"
)

func Add(interactionState *state.InteractionState) {
	opts := interactionState.Interaction.ApplicationCommandData().Options
	num1 := opts[0].IntValue()
	num2 := opts[1].IntValue()
	sum := num1 + num2
	sumString := fmt.Sprintf("%d + %d = %d", num1, num2, sum)

	interactionState.Args["num1"] = strconv.FormatInt(num1, 10)
	interactionState.Args["num2"] = strconv.FormatInt(num2, 10)

	interactionState.Reply(sumString)
}
