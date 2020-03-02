package twitchbot

import (
	"fmt"

	"github.com/gempir/go-twitch-irc"

	"prestrafe-bot/config"
	"prestrafe-bot/gsi"
)

type botChannel struct {
	name        string
	commands    []ChatCommand
	channelSink ChatMessageSink
}

func newChannel(client *twitch.Client, gsiClient gsi.Client, config *config.ChannelConfig) *botChannel {
	commands := []ChatCommand{
		// Troll commands
		NewGlobalCheckCommand().
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("globalcheck")).
			Build(),

		// Map information commands
		NewMapCommand(gsiClient).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("map")).
			Build(),
		NewTierCommand(gsiClient).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("tier")).
			Build(),

		// Player commands
		NewRankCommand(gsiClient).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("rank")).
			Build(),

		// Record time commands
		NewWRCommand(gsiClient).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("wr")).
			Build(),
		NewBWRCommand(gsiClient).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("wr")).
			Build(),
		NewPBCommand(gsiClient).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("pb")).
			Build(),
		NewBPBCommand(gsiClient).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("pb")).
			Build(),

		// Jump Stat commands
		NewJumpStatCommand(gsiClient, "bh", "bhop", "Bunnyhop", 400).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
		NewJumpStatCommand(gsiClient, "dh", "drophop", "Drop Bunnyhop", 400).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
		NewJumpStatCommand(gsiClient, "laj", "ladderjump", "Ladder Jump", 400).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
		NewJumpStatCommand(gsiClient, "lj", "longjump", "Long Jump", 300).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
		NewJumpStatCommand(gsiClient, "mbh", "multibhop", "Multi Bunnyhop", 400).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
		NewJumpStatCommand(gsiClient, "wj", "weirdjump", "Weird Jump", 400).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
	}

	commands = append(commands, NewHelpCommand(commands).
		WithConfig(config.GetCommandConfig("*")).
		WithConfig(config.GetCommandConfig("help")).
		Build(),
	)

	channelName := config.Name
	return &botChannel{
		channelName,
		commands,
		func(format string, a ...interface{}) {
			client.Say(channelName, fmt.Sprintf(format, a...))
		},
	}
}

func (c *botChannel) handle(user *twitch.User, message *twitch.Message) {
	for _, command := range c.commands {
		if command.TryHandle(user, message, c.channelSink) {
			return
		}
	}
}
