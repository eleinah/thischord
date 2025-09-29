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
	switch a.Key {
	case "error":
		return tint.Attr(ANSIBrightRed, a)
	case "description":
		return tint.Attr(ANSIGray, a)
	case "name":
		return tint.Attr(ANSIGray, a)
	case "options":
		return tint.Attr(ANSIGray, a)
	case "username":
		return tint.Attr(ANSIBrightBlue, a)
	case "command":
		return tint.Attr(ANSIBrightBlue, a)
	case "args":
		return tint.Attr(ANSIBrightBlue, a)
	case "version":
		return tint.Attr(ANSIBrightGreen, a)
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
