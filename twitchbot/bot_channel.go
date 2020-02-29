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

func newChannel(client *twitch.Client, config *config.ChannelConfig) *botChannel {
	// TODO Wire actual GSI config here
	gsiClient := gsi.NewClient("localhost", 8337, config.GsiToken)

	commands := []ChatCommand{
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
		NewGlobalCheckCommand().
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("globalcheck")).
			Build(),
		NewMapCommand(gsiClient).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("map")).
			Build(),
		NewTierCommand(gsiClient).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("tier")).
			Build(),
		NewJumpStatCommand(gsiClient, "bh", "bhop", "Bunnyhop").
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
		NewJumpStatCommand(gsiClient, "bh", "drophop", "Drop Bunnyhop").
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
		NewJumpStatCommand(gsiClient, "laj", "ladderjump", "Ladder Jump").
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
		NewJumpStatCommand(gsiClient, "laj", "ladderjump", "Ladder Jump").
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
		NewJumpStatCommand(gsiClient, "lj", "longjump", "Long Jump").
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
		NewJumpStatCommand(gsiClient, "mbh", "multibhop", "Multi Bunnyhop").
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
		NewJumpStatCommand(gsiClient, "wj", "weirdjump", "Weird Jump").
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("jumpstat")).
			Build(),
	}

	return &botChannel{
		config.Name,
		commands,
		func(format string, a ...interface{}) {
			client.Say(config.Name, fmt.Sprintf(format, a...))
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
