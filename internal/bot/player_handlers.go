package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (b *Bot) onPlayerPause(player disgolink.Player, e lavalink.PlayerPauseEvent) {
	slog.Info("player paused", "event", e)
}

func (b *Bot) onPlayerResume(player disgolink.Player, e lavalink.PlayerResumeEvent) {
	slog.Info("player resumed", "event", e)
}

func (b *Bot) onTrackStart(player disgolink.Player, e lavalink.TrackStartEvent) {
	slog.Info("track started", "event", e)
	info := fmt.Sprintf("Now playing: [`%s`](<%s>)", e.Track.Info.Title, *e.Track.Info.URI)
	b.Client.Rest().CreateMessage(DefaultChannel, discord.NewMessageCreateBuilder().
		SetContent(info).
		Build())
}

func (b *Bot) onTrackEnd(player disgolink.Player, e lavalink.TrackEndEvent) {
	if !e.Reason.MayStartNext() {
		return
	}

	var (
		nextTrack lavalink.Track
		ok        bool
	)
	switch b.Queue.Type {
	case QueueTypeNormal:
		nextTrack, ok = b.Queue.Next()

	case QueueTypeRepeatTrack:
		nextTrack = e.Track

	case QueueTypeRepeatQueue:
		b.Queue.Add(e.Track)
		nextTrack, ok = b.Queue.Next()
	}

	if !ok {
		b.Client.Rest().CreateMessage(DefaultChannel, discord.NewMessageCreateBuilder().
			SetContent("Queue concluded.").
			Build())
		return
	}
	if err := player.Update(context.TODO(), lavalink.WithTrack(nextTrack)); err != nil {
		slog.Error("failed to play next track", "error", err)
	}
}

func (b *Bot) onTrackException(player disgolink.Player, e lavalink.TrackExceptionEvent) {
	slog.Info("track exception", "event", e)
}

func (b *Bot) onTrackStuck(player disgolink.Player, e lavalink.TrackStuckEvent) {
	slog.Info("track stuck", "event", e)
}

func (b *Bot) onWebSocketClosed(player disgolink.Player, e lavalink.WebSocketClosedEvent) {
	slog.Info("websocket closed", "event", e)
}

func (b *Bot) onUnknownEvent(p disgolink.Player, e lavalink.UnknownEvent) {
	slog.Info("unknown event", "event", e.Type(), "data", string(e.Data))
}
