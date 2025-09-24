package state

import (
	"log/slog"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type InteractionState struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	User        *discordgo.User
	GuildID     string
	ChannelID   string
	CommandName string
	Args        map[string]string
	Responded   bool
}

func (s *InteractionState) ArgumentstoString() string {
	output := " "
	for k, v := range s.Args {
		if v != "" {
			output += k + "=" + v + " "
		}
	}
	return output
}

func (s *InteractionState) Reply(content string) {
	err := s.Session.InteractionRespond(s.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		slog.Error("Error sending reply", "error", err)
	}
	s.Responded = true
}

func (s *InteractionState) InvisibleReply(content string) {
	err := s.Session.InteractionRespond(s.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   1 << 6,
		},
	})
	if err != nil {
		slog.Error("Error sending reply", "error", err)
	}
	s.Responded = true
}

func (s *InteractionState) Defer() {
	data := &discordgo.InteractionResponseData{}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: data,
	}

	err := s.Session.InteractionRespond(s.Interaction.Interaction, resp)
	if err != nil {
		slog.Error("Error deferring interaction", "error", err)
	}
	s.Responded = true
}

func (s *InteractionState) EditDeferred(content string) {
	hook := &discordgo.WebhookEdit{
		Content: &content,
	}

	if _, err := s.Session.InteractionResponseEdit(s.Interaction.Interaction, hook); err != nil {
		slog.Error("Error editing deferred response", "error", err)
	}
}

func NewInteractionState(s *discordgo.Session, i *discordgo.InteractionCreate) *InteractionState {
	interactionState := &InteractionState{
		Session:     s,
		Interaction: i,
		User:        i.User,
		GuildID:     i.GuildID,
		ChannelID:   i.ChannelID,
		CommandName: i.ApplicationCommandData().Name,
		Args:        make(map[string]string),
		Responded:   false,
	}

	if data := i.ApplicationCommandData(); data.Name != "" {
		for _, option := range data.Options {
			switch option.Type {
			case discordgo.ApplicationCommandOptionInteger:
				interactionState.Args[option.Name] = strconv.FormatInt(option.IntValue(), 10)
			}
		}
	}
	if interactionState.User == nil && i.Member != nil {
		interactionState.User = i.Member.User
	}

	return interactionState
}
