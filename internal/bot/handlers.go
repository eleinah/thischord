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

func ytSearch(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error {
	Defer(event)

	query := data.String("query")
	maxResults := 5

	if len(data.Options) > 1 {
		maxResults = data.Int("limit")
		if maxResults != 0 && maxResults < 0 || maxResults > 25 {
			slog.Debug("limit must be between 0 and 25")
			return EditDeferred("Limit must be between 0 and 25", event)
		}
	}

	search := fmt.Sprintf("ytsearch%d:%s", maxResults, query)

	out, err := ytdlp.New().
		FlatPlaylist().
		Print("title,urls").
		Run(context.Background(), search)

	if err != nil {
		slog.Error("error running ytdlp for query", "error", err)
		return EditDeferred("Error running ytdlp for query", event)
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
		return EditDeferred("Reply too long, consider using a smaller limit.", event)
	}
	return EditDeferred(reply, event)
}
