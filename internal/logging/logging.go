package logging

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func SetDefaultLogger() {
	h := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
	})

	slog.SetDefault(slog.New(h))
}

func FatalLog(msg string, err error) {
	if err != nil {
		slog.Error("FATAL! " + msg + ":: " + err.Error())
		os.Exit(1)
	}
	slog.Error(msg)
	os.Exit(1)
}
