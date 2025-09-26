package handlers

import (
	"fmt"

	"github.com/eleinah/thischord/internal/commands"
	"github.com/eleinah/thischord/internal/state"
)

func commandSelector(interactionState *state.InteractionState) {
	if state.DisabledCommands[interactionState.CommandName] {
		replyString := fmt.Sprintf("'%s' is disabled.", interactionState.CommandName)
		interactionState.InvisibleReply(replyString)
		return
	}

	switch interactionState.CommandName {
	case "join":
		commands.Join(interactionState)
	case "leave":
		commands.Leave(interactionState)
	case "ytsearch":
		commands.YTSearch(interactionState)
	case "ytplay":
		commands.YTPlay(interactionState)
	}
}
