package bot

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
)

func getCommandOptions(cmdIface discord.ApplicationCommandCreate) string {
	body := ""

	switch cmd := cmdIface.(type) {
	case discord.SlashCommandCreate:
		for i, option := range cmd.Options {
			if i < len(cmd.Options)-1 {
				body += fmt.Sprintf("%s,", option.OptionName())
			} else {
				body += fmt.Sprintf("%s", option.OptionName())
			}
		}
	}

	strOptions := fmt.Sprintf("[%s]", body)

	return strOptions
}

func getCommandArgs(cmdName string, data discord.SlashCommandInteractionData) string {
	args := strings.Builder{}

	for _, command := range commands {
		if ok := func(iface discord.ApplicationCommandCreate) bool {
			switch cmd := iface.(type) {
			case discord.SlashCommandCreate:
				if cmd.Name == cmdName {
					for _, opt := range cmd.Options {
						optName := opt.OptionName()
						optType := opt.Type()
						switch optType {
						case discord.ApplicationCommandOptionTypeBool:
							args.WriteString(fmt.Sprintf("'%s'='%v'", optName, data.Bool(optName)))
							args.WriteString(" | ")
						case discord.ApplicationCommandOptionTypeFloat:
							args.WriteString(fmt.Sprintf("'%s'='%v'", optName, data.Float(optName)))
							args.WriteString(" | ")
						case discord.ApplicationCommandOptionTypeInt:
							args.WriteString(fmt.Sprintf("'%s'='%v'", optName, data.Int(optName)))
							args.WriteString(" | ")
						case discord.ApplicationCommandOptionTypeString:
							args.WriteString(fmt.Sprintf("'%s'='%v'", optName, data.String(optName)))
							args.WriteString(" | ")
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
	return strings.TrimSpace(args.String())
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
