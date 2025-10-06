package bot

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// TODO: Finish implementing all queue based commands

func (b *Bot) shuffle(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	// queue := b.Queues.Get(*e.GuildID())
	// if queue == nil {
	// 	return Reply("No player found", e)
	// }

	// queue.Shuffle()
	// return Reply("Queue shuffled", e)
	return InvisibleReply("Not yet implemented", e)
}

func (b *Bot) skip(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	// player := b.Lavalink.ExistingPlayer(*e.GuildID())
	// queue := b.Queues.Get(*e.GuildID())
	// if player == nil || queue == nil {
	// 	return Reply("No player found", e)
	// }

	// amount, ok := d.OptInt("amount")
	// if !ok {
	// 	amount = 1
	// }

	// track, ok := queue.Skip(amount)
	// if !ok {
	// 	return Reply("No tracks in queue", e)
	// }

	// if err := player.Update(context.TODO(), lavalink.WithTrack(track)); err != nil {
	// 	return InvisibleReply(fmt.Sprintf("Error while skipping track: `%s`", err), e)
	// }

	// return Reply("Skipped track", e)
	return InvisibleReply("Not yet implemented", e)
}

func (b *Bot) queueType(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	// queue := b.Queues.Get(*e.GuildID())
	// if queue == nil {
	// 	return Reply("No player found", e)
	// }

	// queue.Type = QueueType(d.String("type"))
	// return Reply(fmt.Sprintf("Queue type set to `%s`", queue.Type), e)
	return InvisibleReply("Not yet implemented", e)
}

func (b *Bot) clearQueue(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	// queue := b.Queues.Get(*e.GuildID())
	// if queue == nil {
	// 	return Reply("No player found", e)
	// }

	// queue.Clear()
	// return Reply("Queue cleared", e)
	return InvisibleReply("Not yet implemented", e)
}

func (b *Bot) queue(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	// queue := b.Queues.Get(*e.GuildID())
	// if queue == nil {
	// 	return Reply("No player found", e)
	// }

	// if len(queue.Tracks) == 0 {
	// 	return Reply("No tracks in queue", e)
	// }

	// var tracks string
	// for i, track := range queue.Tracks {
	// 	tracks += fmt.Sprintf("%d. [`%s`](<%s>)\n", i+1, track.Info.Title, *track.Info.URI)
	// }

	// return Reply(fmt.Sprintf("Queue `%s`:\n%s", queue.Type, tracks), e)
	return InvisibleReply("Not yet implemented", e)
}
