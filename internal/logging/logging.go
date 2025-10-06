package logging

import (
	"errors"
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
)

func colorizeAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case "error":
		return tint.Attr(ANSIRed, a)
	case "err":
		return tint.Attr(ANSIRed, a)
	case "description":
		return tint.Attr(ANSIGray, a)
	case "name":
		return tint.Attr(ANSIGray, a)
	case "options":
		return tint.Attr(ANSIGray, a)
	case "username":
		return tint.Attr(ANSIBlue, a)
	case "command":
		return tint.Attr(ANSIBlue, a)
	case "args":
		return tint.Attr(ANSIBlue, a)
	case "arg":
		return tint.Attr(ANSIRed, a)
	case "stack":
		return tint.Attr(ANSIMagenta, a)
	case "version":
		return tint.Attr(ANSIGreen, a)
	case "node_version":
		return tint.Attr(ANSIGreen, a)
	case "node_session_id":
		return tint.Attr(ANSIGreen, a)
	case "event":
		return tint.Attr(ANSIYellow, a)
	}
	return a
}

func SetDefaultLogger(mode string) {
	opts := &tint.Options{
		TimeFormat:  time.Kitchen,
		ReplaceAttr: colorizeAttr,
	}

	switch mode {
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	default:
		FatalLog("unknown mode", errors.New("invalid logger level"))
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
