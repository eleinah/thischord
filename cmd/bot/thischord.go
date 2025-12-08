package main

import (
	"log/slog"

	"github.com/disgoorg/disgo"
	"github.com/eleinah/thischord/internal/bot"
	"github.com/eleinah/thischord/internal/logging"
)

func main() {
	logging.SetDefaultLogger("info")
	slog.Info("starting ThisChord bot")
	slog.Info("disgo version", "version", disgo.Version)
	bot.Run()
}
