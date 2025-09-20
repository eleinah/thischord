package bot

import (
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/eleinah/thischord/internal/commands"
	"github.com/eleinah/thischord/internal/handlers"
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

	if _, err := exec.LookPath("yt-dlp"); err != nil {
		logging.FatalLog("yt-dlp not found. Please install it: https://github.com/yt-dlp/yt-dlp/wiki/Installation\n", err)
	}

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		logging.FatalLog("ffmpeg not found. Please install it with your package manager or visit https://www.ffmpeg.org/download.html\n", err)
	}

	if _, err := exec.LookPath("ffprobe"); err != nil {
		slog.Warn("ffprobe not found. This may cause problems with some commands.")
	}

	disabled := os.Getenv("DISABLED_COMMANDS")
	for _, command := range strings.Split(disabled, ",") {
		command = strings.TrimSpace(command)
		if command != "" {
			state.DisabledCommands[command] = true
		}
	}
}

func Run() {
	setup()

	dg, err := discordgo.New("Bot " + state.Token)
	if err != nil {
		logging.FatalLog("Error creating Discord session", err)
	}

	dg.AddHandler(handlers.HandleInteractionCreate)

	err = dg.Open()
	if err != nil {
		logging.FatalLog("Error opening connection", err)
	}
	defer func(dg *discordgo.Session) {
		err := dg.Close()
		if err != nil {
			logging.FatalLog("Error closing connection", err)
		}
	}(dg)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	commands.SetupSlashCommands(dg)
	slog.Info("Bot is now running. Press CTRL-C to exit.")
	<-exit
}
