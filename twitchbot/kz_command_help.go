package twitchbot

import (
	"sort"
	"strings"
)

func NewHelpCommand(commands []ChatCommand) ChatCommandBuilder {
	return NewChatCommandBuilder("prestrafe").
		WithAlias("prestrafebot", "kz", "kzhelp", "helpkz", "kzcommands").
		WithHandler(func(ctx CommandContext) (msg string, err error) {
			commandStrings := make([]string, 0)

			for _, command := range commands {
				commandStrings = append(commandStrings, commandPrefix+command.Name())
			}

			sort.Strings(commandStrings)
			return "Available commands: " + strings.Join(commandStrings, ", "), nil
		})
}
