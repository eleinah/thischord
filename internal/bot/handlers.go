package bot

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/lrstanley/go-ytdlp"
)

const (
	botMuted    bool = false
	botDeafened bool = true
)

func (b *Bot) join(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	var voiceState discord.VoiceState
	var ok bool

	if voiceState, ok = b.isUserInVoice(e); !ok {
		return InvisibleReply("You must be in a voice channel to use this command.", e)
	}

	if err := b.Client.UpdateVoiceState(context.TODO(), *e.GuildID(), voiceState.ChannelID, botMuted, botDeafened); err != nil {
		slog.Error("failed to join voice channel", "error", err)
		return Reply("Failed to join the voice channel.", e)
	}

	return Reply("Joined voice channel.", e)
}

func (b *Bot) leave(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	if _, ok := b.isUserInVoice(e); !ok {
		return InvisibleReply("You must be in a voice channel to use this command.", e)
	}

	if err := b.Client.UpdateVoiceState(context.TODO(), *e.GuildID(), nil, botMuted, botDeafened); err != nil {
		slog.Error("failed to join voice channel", "error", err)
		return Reply("Failed to join the voice channel.", e)
	}

	return Reply("Left voice channel.", e)
}

func ytSearch(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	Defer(e)

	query := d.String("query")
	maxResults := 5

	if len(d.Options) > 1 {
		maxResults = d.Int("limit")
		if maxResults != 0 && maxResults < 0 || maxResults > 25 {
			slog.Debug("limit must be between 0 and 25")
			return EditDeferred("Limit must be between 0 and 25", e)
		}
	}

	search := fmt.Sprintf("ytsearch%d:%s", maxResults, query)

	out, err := ytdlp.New().
		FlatPlaylist().
		Print("title,urls").
		Run(context.Background(), search)

	if err != nil {
		slog.Error("error running ytdlp for query", "error", err)
		return EditDeferred("Error running ytdlp for query", e)
	}

	rawResults := strings.Split(out.Stdout, "\n")

	results := make(map[string]string)

	for i := range rawResults {
		if i%2 == 0 {
			results[rawResults[i]] = rawResults[i+1]
		}
	}

	reply := ""
	counter := 1

	for k, v := range results {
		reply += fmt.Sprintf("%d. ***%s***: <%s>\n", counter, k, v)
		counter++
	}

	if len(reply) > 1999 {
		slog.Debug("reply too long, consider using a smaller limit")
		return EditDeferred("Reply too long, consider using a smaller limit.", e)
	}
	return EditDeferred(reply, e)
}
