package bot

import (
	"errors"
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

func Reply(content string, e *events.ApplicationCommandInteractionCreate) error {
	if err := e.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent(content).
		Build(),
	); err != nil {
		slog.Error("error sending response", "error", err)
		return errors.New("error sending response")
	}
	return nil
}

func InvisibleReply(content string, e *events.ApplicationCommandInteractionCreate) error {
	if err := e.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent(content).
		SetEphemeral(true).
		Build(),
	); err != nil {
		slog.Error("error sending response", "error", err)
		return errors.New("error sending response")
	}
	return nil
}

func Defer(e *events.ApplicationCommandInteractionCreate) error {
	if err := e.DeferCreateMessage(false); err != nil {
		slog.Error("error deferring interaction", "error", err)
		return errors.New("error deferring interaction")
	}
	return nil
}

func EditDeferred(content string, e *events.ApplicationCommandInteractionCreate) error {
	if _, err := e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.MessageUpdate{
		Content: &content,
	}); err != nil {
		slog.Error("error editing deferred interaction", "error", err)
		return errors.New("error editing deferred interaction")
	}
	return nil
}
