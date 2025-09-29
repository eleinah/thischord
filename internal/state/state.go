package state

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var (
	Token            string
	DisabledCommands = make(map[string]bool)
	GuildID          snowflake.ID
)

func Reply(content string, e *events.ApplicationCommandInteractionCreate) {
	if err := e.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent(content).
		Build(),
	); err != nil {
		slog.Error("Error sending response", "error", err)
	}
}

func InvisibleReply(content string, e *events.ApplicationCommandInteractionCreate) {
	if err := e.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent(content).
		SetEphemeral(true).
		Build(),
	); err != nil {
		slog.Error("Error sending response", "error", err)
	}
}

func Defer(e *events.ApplicationCommandInteractionCreate) {
	if err := e.DeferCreateMessage(false); err != nil {
		slog.Error("Error deferring interaction", "error", err)
	}
}

func EditDeferred(content string, e *events.ApplicationCommandInteractionCreate) {
	if _, err := e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.MessageUpdate{
		Content: &content,
	}); err != nil {
		slog.Error("Error editing deferred interaction", "error", err)
	}
}
