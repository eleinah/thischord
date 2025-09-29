package bot

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/eleinah/thischord/internal/commands"
	"github.com/eleinah/thischord/internal/logging"
	"github.com/eleinah/thischord/internal/state"
	"github.com/joho/godotenv"
)

func setup() {
	if err := godotenv.Load("../../.env"); err != nil {
		logging.FatalLog("Error loading .env file", err)
	}

	state.Token = os.Getenv("DISCORD_BOT_TOKEN")
	if state.Token == "" {
		logging.FatalLog("No Discord bot token found in .env file", nil)
	}

	state.GuildID = snowflake.GetEnv("DISCORD_GUILD_ID")

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

	client, err := disgo.New(state.Token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentsGuild)),
		bot.WithEventListenerFunc(commands.CommandListener),
	)
	if err != nil {
		logging.FatalLog("Error creating discord client", err)
		return
	}

	defer client.Close(context.TODO())

	commands.SetupSlashCommands(client)

	if err = client.OpenGateway(context.TODO()); err != nil {
		logging.FatalLog("Error opening gateway", err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("Bot is now running. Press CTRL-C to exit.")
	<-exit
}
