package commands

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/eleinah/thischord/internal/state"
	"github.com/lrstanley/go-ytdlp"
)

func YTSearch(d *discord.SlashCommandInteractionData, e *events.ApplicationCommandInteractionCreate) {
	state.Defer(e)

	query := d.String("query")
	maxResults := 5

	if len(d.Options) > 1 {
		maxResults = d.Int("limit")
		if maxResults != 0 && maxResults < 0 || maxResults > 25 {
			state.EditDeferred("Limit must be between 0 and 25", e)
			slog.Debug("Limit must be between 0 and 25")
			return
		}
	}

	search := fmt.Sprintf("ytsearch%d:%s", maxResults, query)

	out, err := ytdlp.New().
		FlatPlaylist().
		Print("title,urls").
		Run(context.Background(), search)

	if err != nil {
		slog.Error("Error running ytdlp for query", "error", err)
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
		slog.Debug("Reply too long, consider using a smaller limit.")
		state.EditDeferred("Reply too long, consider using a smaller limit.", e)
	}
	state.EditDeferred(reply, e)
}
