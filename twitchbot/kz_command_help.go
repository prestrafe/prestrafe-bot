package twitchbot

import (
	"sort"
	"strings"
)

func NewHelpCommand(commands []ChatCommand) ChatCommandBuilder {
	return NewChatCommandBuilder("kz").
		WithAlias("kzhelp", "helpkz", "kzcommands").
		WithHandler(func(ctx CommandContext) (msg string, err error) {
			commandStrings := make([]string, 0)

			for _, command := range commands {
				if command.Enabled() {
					commandHelp := commandPrefix + command.Name()

					if command.SubOnly() {
						commandHelp += " (sub only)"
					}

					commandStrings = append(commandStrings, commandHelp)
				}
			}

			sort.Strings(commandStrings)
			return "Available commands: " + strings.Join(commandStrings, ", "), nil
		})
}
