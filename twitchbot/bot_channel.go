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
		NewPBCommand(gsiClient).
			WithConfig(config.GetCommandConfig("*")).
			WithConfig(config.GetCommandConfig("pb")).
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
