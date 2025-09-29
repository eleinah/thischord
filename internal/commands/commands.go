package commands

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/eleinah/thischord/internal/logging"
	"github.com/eleinah/thischord/internal/state"
)

func CommandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()

	switch data.CommandName() {
	case "ytsearch":
		YTSearch(&data, event)
	}
}

func SetupSlashCommands(client bot.Client) {
	slog.Info("Setting up slash commands...")

	commands := []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "join",
			Description: "join a voice channel",
		},

		discord.SlashCommandCreate{
			Name:        "leave",
			Description: "leave a voice channel",
		},

		discord.SlashCommandCreate{
			Name:        "ytsearch",
			Description: "query YouTube for a song",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "query",
					Description: "song name",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "limit",
					Description: "number of results to return (default: 5, max: 50)",
					Required:    false,
				},
			},
		},

		discord.SlashCommandCreate{
			Name:        "ytplay",
			Description: "play a YouTube video",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "url",
					Description: "youtube url eg. https://www.youtube.com/watch?=???????????",
					Required:    true,
				},
			},
		},
	}

	if _, err := client.Rest().SetGuildCommands(client.ApplicationID(), state.GuildID, commands); err != nil {
		logging.FatalLog("Error registering commands", err)
	}
}
