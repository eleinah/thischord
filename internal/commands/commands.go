package commands

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/eleinah/thischord/internal/logging"
	"github.com/eleinah/thischord/internal/state"
)

func SetupSlashCommands(s *discordgo.Session) {
	slog.Info("Setting up slash commands...")

	commands := []*discordgo.ApplicationCommand{
		{Name: "ping", Description: "pong!"},
		{Name: "greet", Description: "greet the ol' fashion way"},
		{Name: "add", Description: "add two numbers together",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "first",
					Description: "first number",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "second",
					Description: "second number",
					Required:    true,
				}},
		}}

	existingCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		logging.FatalLog("Error getting existing commands", err)
	}

	for _, command := range commands {
		if state.DisabledCommands[command.Name] {
			slog.Warn("Skipping disabled command: " + command.Name)
			for _, existingCommand := range existingCommands {
				if existingCommand.Name == command.Name {
					err := s.ApplicationCommandDelete(s.State.User.ID, existingCommand.GuildID, existingCommand.ID)
					if err != nil {
						slog.Error("Error deleting command:", err)
					}
					slog.Info("Deleted command: " + command.Name)
				}
			}
			continue
		}
		found := false
		for _, existingCommand := range existingCommands {
			if existingCommand.Name == command.Name {
				found = true
				slog.Info("Retrieved command: " + command.Name)
				break
			}
		}
		if !found {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
			if err != nil {
				logging.FatalLog("Error creating command", err)
			} else {
				slog.Info("Created command: " + command.Name)
			}
		}
	}

	slog.Info("Finished setting up slash commands")
}
