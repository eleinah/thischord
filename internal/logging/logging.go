package logging

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

const (
	ANSIBlack uint8 = iota
	ANSIRed
	ANSIGreen
	ANSIYellow
	ANSIBlue
	ANSIMagenta
	ANSICyan
	ANSIWhite
	ANSIGray
	ANSIBrightRed
	ANSIBrightGreen
	ANSIBrightYellow
	ANSIBrightBlue
	ANSIBrightMagenta
	ANSIBrightCyan
	ANSIBrightWhite
)

func colorizeAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == "error" {
		return tint.Attr(ANSIBrightRed, a)
	}
	if a.Key == "description" {
		return tint.Attr(ANSIGray, a)
	}
	if a.Key == "name" {
		return tint.Attr(ANSIGray, a)
	}
	if a.Key == "options" {
		return tint.Attr(ANSIGray, a)
	}
	return a
}

func SetDefaultLogger() {
	opts := &tint.Options{
		Level:       slog.LevelInfo,
		TimeFormat:  time.Kitchen,
		ReplaceAttr: colorizeAttr,
	}

	h := tint.NewHandler(os.Stderr, opts)

	slog.SetDefault(slog.New(h))
}

func FatalLog(msg string, err error) {
	if err != nil {
		slog.Error("FATAL! "+msg+":: ", "error", err.Error())
		os.Exit(1)
	}
	slog.Error("FATAL! " + msg)
	os.Exit(1)
}
