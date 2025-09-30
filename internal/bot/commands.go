package bot

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/eleinah/thischord/internal/logging"
)

var commands = []discord.ApplicationCommandCreate{
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
