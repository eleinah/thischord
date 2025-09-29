package commands

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/eleinah/thischord/internal/state"
)

func DeferMsg(d *discord.SlashCommandInteractionData, e *events.ApplicationCommandInteractionCreate) {
	state.Defer(e)

	time.Sleep(5 * time.Second)

	state.EditDeferred("I have awoken after 5 seconds!", e)
}
