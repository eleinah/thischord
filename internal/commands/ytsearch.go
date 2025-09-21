package commands

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/eleinah/thischord/internal/state"
	"github.com/lrstanley/go-ytdlp"
)

func YTSearch(interactionState *state.InteractionState) {
	opts := interactionState.Interaction.ApplicationCommandData().Options

	query := opts[0].StringValue()
	maxResults := int64(5)

	if len(opts) > 1 {
		maxResults = opts[1].IntValue()
		if int(maxResults) != 0 && int(maxResults) < 0 || int(maxResults) > 50 {
			interactionState.Reply("Limit must be between 0 and 50")
		}
	}

	search := fmt.Sprintf("ytsearch%d:%s", maxResults, query)

	out, err := ytdlp.New().
		FlatPlaylist().
		Print("title,urls").
		Run(context.Background(), search)

	if err != nil {
		slog.Error("Error running ytdlp for query:", err)
	}

	results := strings.Split(out.Stdout, "\n")

	reply := ""
	counter := 0

	for i, result := range results {
		if i%2 == 0 {
			reply += fmt.Sprintf("[%d] Title: %s\n", counter+1, result)
			counter++
		} else {
			reply += fmt.Sprintf("URL: <%s>\n\n", result)
		}
	}

	interactionState.Reply(reply)
}
