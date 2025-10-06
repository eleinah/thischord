package bot

import (
	"context"
	"log/slog"

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
}

func (b *Bot) onTrackEnd(player disgolink.Player, e lavalink.TrackEndEvent) {
	if !e.Reason.MayStartNext() {
		return
	}

	queue := b.Queues.Get(e.GuildID())
	var (
		nextTrack lavalink.Track
		ok        bool
	)
	switch queue.Type {
	case QueueTypeNormal:
		nextTrack, ok = queue.Next()

	case QueueTypeRepeatTrack:
		nextTrack = e.Track

	case QueueTypeRepeatQueue:
		queue.Add(e.Track)
		nextTrack, ok = queue.Next()
	}

	if !ok {
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
