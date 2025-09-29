package bot

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/eleinah/thischord/internal/logging"
	"github.com/joho/godotenv"
)

func newBot() *Bot {
	return &Bot{}
}

type Bot struct {
	Client   bot.Client
	Handlers map[string]func(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error
}

func (b *Bot) onApplicationCommand(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()

	slog.Info("interaction received", "username", event.User().Username, "command", data.CommandName(), "args", getCommandArgs(data.CommandName(), data))

	handler, ok := b.Handlers[data.CommandName()]
	if !ok {
		slog.Info("unknown command", "command", data.CommandName())
		return
	}
	if err := handler(event, data); err != nil {
		slog.Error("error handling command", "error", err)
	}
}

func setup() {
	if err := godotenv.Load("../../.env"); err != nil {
		logging.FatalLog("error loading .env file", err)
	}

	Token = os.Getenv("DISCORD_BOT_TOKEN")
	if Token == "" {
		logging.FatalLog("no Discord bot token found in .env file", nil)
	}

	GuildID = snowflake.GetEnv("DISCORD_GUILD_ID")

	if _, err := exec.LookPath("yt-dlp"); err != nil {
		logging.FatalLog("yt-dlp not found. Please install it: https://github.com/yt-dlp/yt-dlp/wiki/Installation\n", err)
	}

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		logging.FatalLog("ffmpeg not found. Please install it with your package manager or visit https://www.ffmpeg.org/download.html\n", err)
	}

	if _, err := exec.LookPath("ffprobe"); err != nil {
		slog.Warn("ffprobe not found. This may cause problems with some commands.")
	}
}

func Run() {
	setup()

	b := newBot()

	client, err := disgo.New(Token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentsGuild)),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagVoiceStates)),
		bot.WithEventListenerFunc(b.onApplicationCommand),
	)
	if err != nil {
		logging.FatalLog("error creating discord client", err)
		return
	}
	b.Client = client
	defer client.Close(context.TODO())

	registerCommands(client)
	b.Handlers = map[string]func(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error{
		"ytsearch": ytSearch,
	}
	logHandledCommands(b.Handlers)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.OpenGateway(ctx); err != nil {
		logging.FatalLog("error opening gateway", err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("bot is now running, press CTRL-C to exit")
	<-exit
}
