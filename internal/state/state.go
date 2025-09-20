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
	RawArgs     map[string]any
	Args        map[string]string
	Responded   bool
}

func (s *InteractionState) ArgumentstoString() string {
	var output string
	for k, v := range s.Args {
		if v != "" {
			output += " " + k + "=" + v
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
		slog.Error("Error sending reply: ", err)
	}
	s.Responded = true
}

func (s *InteractionState) standardizeArguments() {
	switch s.CommandName {
	case "add":
		if val, exists := s.RawArgs["num1"]; exists {
			switch v := val.(type) {
			case int64:
				s.Args["num1"] = strconv.Itoa(int(v))
			}
		}
		if val, exists := s.RawArgs["num2"]; exists {
			switch v := val.(type) {
			case int64:
				s.Args["num2"] = strconv.Itoa(int(v))
			}
		}
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
		RawArgs:     make(map[string]any),
		Args:        make(map[string]string),
		Responded:   false,
	}

	if data := i.ApplicationCommandData(); data.Name != "" {
		for _, option := range data.Options {
			interactionState.RawArgs[option.Name] = option.Value
		}
	}
	if interactionState.User == nil && i.Member != nil {
		interactionState.User = i.Member.User
	}

	interactionState.standardizeArguments()
	return interactionState
}
