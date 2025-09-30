package bot

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (b *Bot) isUserInVoice(e *events.ApplicationCommandInteractionCreate) (discord.VoiceState, bool) {
	voiceState, ok := b.Client.Caches().VoiceState(*e.GuildID(), e.User().ID)
	if !ok {
		return discord.VoiceState{}, false
	}

	return voiceState, ok
}

func getCommandOptions(cmdIface discord.ApplicationCommandCreate) string {
	body := strings.Builder{}

	body.WriteString("[")
	switch cmd := cmdIface.(type) {
	case discord.SlashCommandCreate:
		for i, option := range cmd.Options {
			if i < len(cmd.Options)-1 {
				body.WriteString(fmt.Sprintf("%s,", option.OptionName()))
			} else {
				body.WriteString(fmt.Sprintf("%s", option.OptionName()))
			}
		}
	}
	body.WriteString("]")

	return body.String()
}

func getCommandArgs(cmdName string, data discord.SlashCommandInteractionData) string {
	var args strings.Builder

	for _, command := range commands {
		if ok := func(iface discord.ApplicationCommandCreate) bool {
			switch cmd := iface.(type) {
			case discord.SlashCommandCreate:
				if cmd.Name == cmdName {
					for _, opt := range cmd.Options {
						optName := opt.OptionName()
						optType := opt.Type()
						optStr := "%s='%v'"
						switch optType {
						case discord.ApplicationCommandOptionTypeBool:
							args.WriteString(fmt.Sprintf(optStr, optName, data.Bool(optName)))
							args.WriteString(", ")
						case discord.ApplicationCommandOptionTypeFloat:
							args.WriteString(fmt.Sprintf(optStr, optName, data.Float(optName)))
							args.WriteString(", ")
						case discord.ApplicationCommandOptionTypeInt:
							switch optName {
							case "limit":
								if data.Int(optName) == 0 {
									args.WriteString(fmt.Sprintf(optStr, optName, "default"))
									args.WriteString(", ")
								} else {
									args.WriteString(fmt.Sprintf(optStr, optName, data.Int(optName)))
									args.WriteString(", ")
								}
							default:
								args.WriteString(fmt.Sprintf(optStr, optName, data.Int(optName)))
								args.WriteString(", ")
							}
						case discord.ApplicationCommandOptionTypeString:
							args.WriteString(fmt.Sprintf(optStr, optName, data.String(optName)))
							args.WriteString(", ")
						default:
							return false
						}
					}
				}
				return true
			default:
				return false
			}
		}(command); !ok {
			return ""
		}
	}

	strArgs := strings.TrimSpace(args.String())
	if len(strArgs) != 0 {
		return strArgs[:len(strArgs)-1]
	}
	return strArgs
}

func getCommandName(iface discord.ApplicationCommandCreate) string {
	name := ""
	switch cmd := iface.(type) {
	case discord.SlashCommandCreate:
		name = cmd.Name
	}
	return name
}

func getCommandDescription(iface discord.ApplicationCommandCreate) string {
	desc := ""
	switch cmd := iface.(type) {
	case discord.SlashCommandCreate:
		desc = cmd.Description
	}
	return desc
}
