package twitchbot

import (
	"github.com/gempir/go-twitch-irc/v2"
	"strings"
)

// This interface defines the public API of the chat bot. The API is pretty slim, as most of the work is done
// internally.
type ChatBot interface {
	// Joins a channel that is defined inside the passed channel configuration.
	Join(channel string, commands []ChatCommand) ChatBot
	// Starts up the chat bot in the thread that is calling this method. It blocks until an error occurs or the bot is
	// stopped.
	Start() error
	// Stops the bot and disconnects it from the Twitch API.
	Stop() error
}

type chatBot struct {
	channels map[string]botChannel
	client   *twitch.Client
}

// Creates a new chat bot instance.
func NewChatBot(userName, accessToken string) ChatBot {
	return &chatBot{
		make(map[string]botChannel),
		twitch.NewClient(userName, accessToken),
	}
}

func (c chatBot) Join(channel string, commands []ChatCommand) ChatBot {
	botChannel := newChannel(channel, c.client, commands)
	channelName := strings.ToLower(botChannel.name)

	c.client.Join(channelName)
	c.channels[channelName] = *botChannel

	return c
}

func (c chatBot) Start() error {
	c.client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if botChannel, hasChannel := c.channels[message.Channel]; hasChannel {
			botChannel.handle(&message)
		}
	})
	return c.client.Connect()
}

func (c chatBot) Stop() error {
	return c.client.Disconnect()
}
