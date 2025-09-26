package commands

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/eleinah/thischord/internal/helpers"
	"github.com/eleinah/thischord/internal/logging"
	"github.com/eleinah/thischord/internal/state"
)

func SetupSlashCommands(s *discordgo.Session) {
	slog.Info("Setting up slash commands...")

	commands := []*discordgo.ApplicationCommand{
		{Name: "join", Description: "join a voice channel"},
		{Name: "leave", Description: "leave a voice channel"},
		{Name: "ytsearch", Description: "query YouTube for a song",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "song name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "limit",
					Description: "number of results to return (default: 5, max: 50)",
					Required:    false,
				},
			},
		},
		{Name: "ytplay", Description: "play a YouTube video",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "youtube url eg. https://www.youtube.com/watch?v=???????????",
					Required:    true,
				},
			},
		},
	}

	existingCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		logging.FatalLog("Error getting existing commands", err)
	}

	for _, command := range commands {
		strOptions := helpers.ListCommandOptions(command.Options)

		if state.DisabledCommands[command.Name] {
			slog.Warn("...skipping disabled command...", "name", command.Name, "options", strOptions, "description", command.Description)
			for _, existingCommand := range existingCommands {
				if existingCommand.Name == command.Name {
					err := s.ApplicationCommandDelete(s.State.User.ID, existingCommand.GuildID, existingCommand.ID)
					if err != nil {
						slog.Error("Error deleting command:", "error", err)
					}
					slog.Info("...deleted command...", "name", command.Name, "options", strOptions, "description", command.Description)
				}
			}
			continue
		}
		found := false
		for _, existingCommand := range existingCommands {
			if existingCommand.Name == command.Name {
				found = true
				slog.Info("...retrieved command...", "name", command.Name, "options", strOptions, "description", command.Description)
				break
			}
		}
		if !found {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
			if err != nil {
				logging.FatalLog("Error creating command", err)
			} else {
				slog.Info("...created command...", "name", command.Name, "options", strOptions, "description", command.Description)
			}
		}
	}

	slog.Info("Finished setting up slash commands!")
}
