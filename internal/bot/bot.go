package bot

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/eleinah/thischord/internal/logging"
	"github.com/joho/godotenv"
)

func setup() {
	if err := godotenv.Load(".env"); err != nil {
		logging.FatalLog("error loading .env file", err)
	}

	Token = os.Getenv("DISCORD_BOT_TOKEN")
	if Token == "" {
		logging.FatalLog("no Discord bot token found in .env file", nil)
	}

	rawGuildID := os.Getenv("DISCORD_GUILD_ID")
	if rawGuildID == "" {
		logging.FatalLog("no Discord guild ID found in .env file", nil)
	}
	GuildID = snowflake.MustParse(rawGuildID)

	NodeName = os.Getenv("NODE_NAME")
	if NodeName == "" {
		logging.FatalLog("no Lavalink node name found in .env file", nil)
	}

	NodeAddress = os.Getenv("NODE_ADDRESS")
	if NodeAddress == "" {
		logging.FatalLog("no Lavalink node address found in .env file", nil)
	}

	NodePassword = os.Getenv("NODE_PASSWORD")
	if NodePassword == "" {
		logging.FatalLog("no Lavalink node password found in .env file", nil)
	}

	rawNodeSecure := os.Getenv("NODE_SECURE")
	if rawNodeSecure == "" {
		logging.FatalLog("no Lavalink secure setting found in .env file", nil)
	}
	NodeSecure, _ = strconv.ParseBool(rawNodeSecure)

	var ignore bool
	var rawIgnore string = os.Getenv("IGNORE_FAILOVER")
	if rawIgnore != "" {
		ignore, _ = strconv.ParseBool(rawIgnore)
	}
	if rawIgnore == "" || !ignore {
		if _, err := exec.LookPath("yt-dlp"); err != nil {
			logging.FatalLog("yt-dlp not found. Please install it: https://github.com/yt-dlp/yt-dlp/wiki/Installation\n", err)
		}

		if _, err := exec.LookPath("ffmpeg"); err != nil {
			logging.FatalLog("ffmpeg not found. Please install it with your package manager or visit https://www.ffmpeg.org/download.html\n", err)
		}

		if _, err := exec.LookPath("ffprobe"); err != nil {
			slog.Warn("ffprobe not found. This may cause problems with some commands.")
		}
	} else {
		slog.Warn("ignoring failover (yt-dlp & ffmpeg)")
	}
}

func Run() {
	setup()

	b := newBot()

	client, err := disgo.New(Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildVoiceStates),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates),
		),
		bot.WithEventListenerFunc(b.onApplicationCommand),
		bot.WithEventListenerFunc(b.onVoiceStateUpdate),
		bot.WithEventListenerFunc(b.onVoiceServerUpdate),
	)
	if err != nil {
		slog.Error("error starting disgo client", "error", err)
		os.Exit(1)
	}
	b.Client = client

	registerCommands(client)

	b.Lavalink = disgolink.New(client.ApplicationID(),
		disgolink.WithListenerFunc(b.onPlayerPause),
		disgolink.WithListenerFunc(b.onPlayerResume),
		disgolink.WithListenerFunc(b.onTrackStart),
		disgolink.WithListenerFunc(b.onTrackEnd),
		disgolink.WithListenerFunc(b.onTrackException),
		disgolink.WithListenerFunc(b.onTrackStuck),
		disgolink.WithListenerFunc(b.onWebSocketClosed),
		disgolink.WithListenerFunc(b.onUnknownEvent),
	)
	b.Handlers = map[string]func(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error{
		"connect":     b.connect,
		"disconnect":  b.disconnect,
		"ytsearch":    b.ytSearch,
		"play":        b.play,
		"pause":       b.pause,
		"now-playing": b.nowPlaying,
		"stop":        b.stop,
		"queue":       b.queue,      // TODO: finish implementing
		"clear-queue": b.clearQueue, // TODO: finish implementing
		"queue-type":  b.queueType,  // TODO: finish implementing
		"shuffle":     b.shuffle,    // TODO: finish implementing
		"skip":        b.skip,       // TODO: finish implementing
	}
	logHandledCommands(b.Handlers)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.OpenGateway(ctx); err != nil {
		slog.Error("failed to open gateway", "error", err)
		os.Exit(1)
	}
	defer client.Close(context.TODO())

	node, err := b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     NodeName,
		Address:  NodeAddress,
		Password: NodePassword,
		Secure:   NodeSecure,
	})
	if err != nil {
		slog.Error("failed to add node", "error", err)
		os.Exit(1)
	}
	version, err := node.Version(ctx)
	if err != nil {
		slog.Error("failed to get node version", "error", err)
		os.Exit(1)
	}

	slog.Info("lavalink node successfully connected", "node_version", version, "node_session_id", node.SessionID())
	slog.Info("bot is now running, press CTRL-C to exit")
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
}

func newBot() *Bot {
	return &Bot{}
}

type Bot struct {
	Client   bot.Client
	Lavalink disgolink.Client
	Handlers map[string]func(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error
	Queues   *QueueManager
}

func (b *Bot) onApplicationCommand(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	cmdName := data.CommandName()

	slog.Info("interaction received", "username", event.User().Username, "command", cmdName, "args", getCommandArgs(cmdName, data))

	handler, ok := b.Handlers[data.CommandName()]
	if !ok {
		slog.Info("unknown command", "command", data.CommandName())
		return
	}
	if err := handler(event, data); err != nil {
		slog.Error("error handling command", "error", err)
	}
}

func (b *Bot) onVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	if event.VoiceState.UserID != b.Client.ApplicationID() {
		return
	}
	b.Lavalink.OnVoiceStateUpdate(context.TODO(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
	if event.VoiceState.ChannelID == nil {
		b.Queues.Delete(event.VoiceState.GuildID)
	}
}

func (b *Bot) onVoiceServerUpdate(event *events.VoiceServerUpdate) {
	b.Lavalink.OnVoiceServerUpdate(context.TODO(), event.GuildID, event.Token, *event.Endpoint)
}
