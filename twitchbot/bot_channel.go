package twitchbot

import (
	"fmt"
	"gitlab.com/prestrafe/prestrafe-bot/utils"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
)

// Buffer sending messages, to avoid spamming twitch and getting banned :(
var messageQueue = utils.CreateTaskQueue(20, 30)

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
			messageQueue.ScheduleTask(func() {
				client.Say(strings.ToLower(channelName), fmt.Sprintf(format, a...))
			})
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
