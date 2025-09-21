package handlers

import (
	"fmt"

	"github.com/eleinah/thischord/internal/commands"
	"github.com/eleinah/thischord/internal/state"
)

func commandSelector(interactionState *state.InteractionState) {
	if state.DisabledCommands[interactionState.CommandName] {
		replyString := fmt.Sprintf("'%s' is disabled.", interactionState.CommandName)
		interactionState.Reply(replyString)
		return
	}

	switch interactionState.CommandName {
	case "ping":
		commands.Ping(interactionState)
	case "greet":
		commands.Greet(interactionState)
	case "add":
		commands.Add(interactionState)
	case "ytsearch":
		commands.Query(interactionState)
	}
}
