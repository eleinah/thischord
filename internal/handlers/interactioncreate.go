package handlers

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/eleinah/thischord/internal/state"
)

func HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		interactionState := state.NewInteractionState(s, i)
		slog.Info("Interaction received:", "username", interactionState.User.Username, "command", interactionState.CommandName, "args", interactionState.ArgumentstoString())
		commandSelector(interactionState)
	}
	return
}
