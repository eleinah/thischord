package state

var (
	Token            string
	DisabledCommands = make(map[string]bool)
	VoiceConnected   bool
	Repeat           bool
)
