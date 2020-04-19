package twitchbot

import (
	"fmt"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
)

type botChannel struct {
	name        string
	commands    []ChatCommand
	channelSink ChatMessageSink
}

func newChannel(channelName string, client *twitch.Client, commands []ChatCommand) *botChannel {
	return &botChannel{
		channelName,
		commands,
		func(format string, a ...interface{}) {
			client.Say(strings.ToLower(channelName), fmt.Sprintf(format, a...))
		},
	}
}

func (c *botChannel) handle(user *twitch.User, message *twitch.PrivateMessage) {
	for _, command := range c.commands {
		if command.TryHandle(c.name, user, message, c.channelSink) {
			return
		}
	}
}
