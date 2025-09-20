package main

import (
	"github.com/eleinah/thischord/internal/bot"
	"github.com/eleinah/thischord/internal/logging"
)

func init() {
	logging.SetDefaultLogger()
}

func main() {
	bot.Run()
}
