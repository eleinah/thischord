package state

import "time"

func init() {
	StartTime = time.Now()
}

var (
	Token            string
	DisabledCommands = make(map[string]bool)
	StartTime        time.Time
	VoiceConnected   bool
	Repeat           bool
)
