package bot

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/lrstanley/go-ytdlp"
)

const (
	botMuted    bool = false
	botDeafened bool = true
)

func (b *Bot) connect(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	vs, ok := b.Client.Caches().VoiceState(*e.GuildID(), e.User().ID)
	if !ok {
		return InvisibleReply("You must be in a voice channel to use this command.", e)
	}

	if err := b.Client.UpdateVoiceState(context.TODO(), *e.GuildID(), vs.ChannelID, botMuted, botDeafened); err != nil {
		slog.Error("failed to join voice channel", "error", err)
		return Reply("Failed to join the voice channel.", e)
	}

	return Reply("Connected to voice channel", e)
}

func (b *Bot) disconnect(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(*e.GuildID())
	if player == nil {
		return Reply("No player found", e)
	}

	if err := b.Client.UpdateVoiceState(context.TODO(), *e.GuildID(), nil, botMuted, botDeafened); err != nil {
		return Reply(fmt.Sprintf("Error while disconnecting: `%s`", err), e)
	}

	return Reply("Disconnected from voice channel", e)
}

func (b *Bot) ytSearch(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	// TODO: refactor this to be generic and use lavalink types
	DeferReply(e)

	query := d.String("query")
	maxResults := 5

	if len(d.Options) > 1 {
		maxResults = d.Int("limit")
		if maxResults != 0 && maxResults < 0 || maxResults > 25 {
			slog.Debug("limit must be between 0 and 25")
			return EditDeferredReply("Limit must be between 0 and 25", e)
		}
	}

	search := fmt.Sprintf("ytsearch%d:%s", maxResults, query)

	out, err := ytdlp.New().
		FlatPlaylist().
		Print("title,urls").
		Run(context.Background(), search)

	if err != nil {
		slog.Error("error running ytdlp for query", "error", err)
		return EditDeferredReply("Error running ytdlp for query", e)
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
		return EditDeferredReply("Reply too long, consider using a smaller limit.", e)
	}
	return EditDeferredReply(reply, e)
}

func (b *Bot) pause(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(*e.GuildID())
	if player == nil {
		return Reply("No player found", e)
	}

	if err := player.Update(context.TODO(), lavalink.WithPaused(!player.Paused())); err != nil {
		return Reply(fmt.Sprintf("Error while pausing: `%s`", err), e)
	}

	status := "playing"
	if player.Paused() {
		status = "paused"
	}
	return Reply(fmt.Sprintf("Player is now %s", status), e)
}

func (b *Bot) stop(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(*e.GuildID())
	if player == nil {
		return Reply("No player found", e)
	}

	if err := player.Update(context.TODO(), lavalink.WithNullTrack()); err != nil {
		return Reply(fmt.Sprintf("Error while stopping: `%s`", err), e)
	}

	return Reply("Player stopped", e)
}

func (b *Bot) nowPlaying(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(*e.GuildID())
	if player == nil {
		return Reply("No player found", e)
	}

	track := player.Track()
	if track == nil {
		return Reply("No track found", e)
	}

	return Reply(fmt.Sprintf("Now playing: [`%s`](<%s>)\n\n %s / %s", track.Info.Title, *track.Info.URI, formatPosition(player.Position()), formatPosition(track.Info.Length)), e)
}

func (b *Bot) play(e *events.ApplicationCommandInteractionCreate, d discord.SlashCommandInteractionData) error {
	identifier := d.String("identifier")
	if source, ok := d.OptString("source"); ok {
		identifier = lavalink.SearchType(source).Apply(identifier)
	} else if !urlPattern.MatchString(identifier) && !searchPattern.MatchString(identifier) {
		identifier = lavalink.SearchTypeYouTube.Apply(identifier)
	}

	voiceState, ok := b.Client.Caches().VoiceState(*e.GuildID(), e.User().ID)
	if !ok {
		return Reply("You need to be in a voice channel to use this command", e)
	}

	DeferReply(e)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var toPlay *lavalink.Track
	b.Lavalink.BestNode().LoadTracksHandler(ctx, identifier, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			EditDeferredReply(fmt.Sprintf("Loaded track: [`%s`](<%s>)", track.Info.Title, *track.Info.URI), e)
			toPlay = &track
		},
		func(playlist lavalink.Playlist) {
			EditDeferredReply(fmt.Sprintf("Loaded playlist: `%s` with `%d` tracks", playlist.Info.Name, len(playlist.Tracks)), e)
			toPlay = &playlist.Tracks[0]
		},
		func(tracks []lavalink.Track) {
			EditDeferredReply(fmt.Sprintf("Loaded search result: [`%s`](<%s>)", tracks[0].Info.Title, *tracks[0].Info.URI), e)
			toPlay = &tracks[0]
		},
		func() {
			EditDeferredReply(fmt.Sprintf("Nothing found for: `%s`", identifier), e)
		},
		func(err error) {
			EditDeferredReply(fmt.Sprintf("Error while looking up query: `%s`", err), e)
		},
	))
	if toPlay == nil {
		return nil
	}

	if err := b.Client.UpdateVoiceState(context.TODO(), *e.GuildID(), voiceState.ChannelID, botMuted, botDeafened); err != nil {
		return err
	}

	return b.Lavalink.Player(*e.GuildID()).Update(context.TODO(), lavalink.WithTrack(*toPlay))
}
