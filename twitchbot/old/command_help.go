package old

import (
	"sort"
	"strings"

	"github.com/gempir/go-twitch-irc"
)

func (bot *Bot) CreateHelpCommand() CommandHandler {
	return func(user twitch.User, parameters []string) string {
		var commands []string

		for name, command := range bot.commands {
			if command.Enabled {
				commandHelp := " !" + name

				if command.SubOnly {
					commandHelp += " (sub only)"
				}

				commands = append(commands, commandHelp)
			}
		}

		sort.Strings(commands)
		return "Available commands: " + strings.Join(commands, ", ")
	}
}
