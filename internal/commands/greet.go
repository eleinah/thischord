package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/eleinah/thischord/internal/state"
)

func Greet(d *discord.SlashCommandInteractionData, e *events.ApplicationCommandInteractionCreate) {
	state.Reply("Hello, world! (this reply is abstracted now!)", e)
}

func InvisibleGreet(d *discord.SlashCommandInteractionData, e *events.ApplicationCommandInteractionCreate) {
	state.InvisibleReply("Hello, world! (this reply is abstracted and invisible now!)", e)
}
