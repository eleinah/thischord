package bot

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/eleinah/thischord/internal/logging"
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "connect",
		Description: "connect player to voice channel",
	},
	discord.SlashCommandCreate{
		Name:        "disconnect",
		Description: "disconnect the player from voice channel",
	},
	discord.SlashCommandCreate{
		// TODO: refactor this to be generic
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
		Name:        "play",
		Description: "play a song",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "identifier",
				Description: "song URI or query",
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				Name:        "source",
				Description: "source to search",
				Required:    false,
				Choices: []discord.ApplicationCommandOptionChoiceString{
					{
						Name:  "YouTube",
						Value: string(lavalink.SearchTypeYouTube),
					},
				},
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "pause",
		Description: "pause the current song",
	},
	discord.SlashCommandCreate{
		Name:        "now-playing",
		Description: "show the song that's playing",
	},
	discord.SlashCommandCreate{
		Name:        "stop",
		Description: "stop the current song and player",
	},
	discord.SlashCommandCreate{
		Name:        "skip",
		Description: "skip the current song",
	},
	discord.SlashCommandCreate{
		Name:        "shuffle",
		Description: "shuffles the queue",
	},
	discord.SlashCommandCreate{
		Name:        "queue",
		Description: "displays the queue",
	},
	discord.SlashCommandCreate{
		Name:        "clear-queue",
		Description: "clears the queue",
	},
}

func logHandledCommands(cmds map[string]func(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error) {
	for name := range cmds {
		for _, cmd := range commands {
			if name == getCommandName(cmd) {
				strOptions := getCommandOptions(cmd)
				desc := getCommandDescription(cmd)
				slog.Info("retrieved command", "name", name, "options", strOptions, "description", desc)
			}
		}
	}
	slog.Info("sucessfully registered commands")
}

func registerCommands(client bot.Client) {
	slog.Info("registering slash commands")
	if err := handler.SyncCommands(client, commands, []snowflake.ID{GuildID}); err != nil {
		logging.FatalLog("error registering commands", err)
	}
}
