package commands

import (
	"fmt"
	"time"

	"github.com/eleinah/thischord/internal/state"
)

func Ping(interactionState *state.InteractionState) {
	start := time.Now()
	interactionState.Reply(fmt.Sprintf("Pong! (took %v)", time.Since(start)), true)
}
