package bot

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (b *Bot) shuffle(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	if b.Queue == nil {
		return Reply("No player found", e)
	}

	b.Queue.Shuffle()
	return Reply("Queue shuffled", e)
}

func (b *Bot) skip(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(*e.GuildID())
	if player == nil || b.Queue == nil {
		return Reply("No player found", e)
	}

	current := b.Lavalink.ExistingPlayer(*e.GuildID()).Track()

	track, ok := b.Queue.Next()
	if !ok {
		return Reply("No tracks in queue", e)
	}

	if err := player.Update(context.TODO(), lavalink.WithTrack(track)); err != nil {
		return InvisibleReply(fmt.Sprintf("Error while skipping track: `%s`", err), e)
	}

	return Reply(fmt.Sprintf("Skipped track [`%s`](<%s>)", current.Info.Title, *current.Info.URI), e)
}

func (b *Bot) queueType(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	if b.Queue == nil {
		return Reply("No player found", e)
	}

	b.Queue.Type = QueueType(d.String("type"))
	return Reply(fmt.Sprintf("Queue type set to `%s`", b.Queue.Type), e)
}

func (b *Bot) clearQueue(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	if b.Queue == nil {
		return Reply("No player found", e)
	}

	b.Queue.Clear()
	return Reply("Queue cleared", e)
}

func (b *Bot) queue(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	if b.Queue == nil {
		return Reply("No player found", e)
	}

	if len(b.Queue.Tracks) == 0 {
		return Reply("No tracks in queue", e)
	}

	var tracks string
	for i, track := range b.Queue.Tracks {
		tracks += fmt.Sprintf("%d. [`%s`](<%s>)\n", i+1, track.Info.Title, *track.Info.URI)
	}

	return Reply(fmt.Sprintf("Queue `%s`:\n%s", b.Queue.Type, tracks), e)
}
