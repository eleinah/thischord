package bot

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func getCommandOptions(iface discord.ApplicationCommandCreate) string {
	body := strings.Builder{}

	body.WriteString("[")
	switch cmd := iface.(type) {
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

// TODO: clean this up when making search command generic
func getCommandArgs(cmdName string, d discord.SlashCommandInteractionData) string {
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
							args.WriteString(fmt.Sprintf(optStr, optName, d.Bool(optName)))
							args.WriteString(", ")
						case discord.ApplicationCommandOptionTypeFloat:
							args.WriteString(fmt.Sprintf(optStr, optName, d.Float(optName)))
							args.WriteString(", ")
						case discord.ApplicationCommandOptionTypeInt:
							switch optName {
							case "limit":
								if d.Int(optName) == 0 {
									args.WriteString(fmt.Sprintf(optStr, optName, "default"))
									args.WriteString(", ")
								} else {
									args.WriteString(fmt.Sprintf(optStr, optName, d.Int(optName)))
									args.WriteString(", ")
								}
							default:
								args.WriteString(fmt.Sprintf(optStr, optName, d.Int(optName)))
								args.WriteString(", ")
							}
						case discord.ApplicationCommandOptionTypeString:
							args.WriteString(fmt.Sprintf(optStr, optName, d.String(optName)))
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

func formatPosition(position lavalink.Duration) string {
	if position == 0 {
		return "0:00"
	}
	return fmt.Sprintf("%d:%02d", position.Minutes(), position.SecondsPart())
}
