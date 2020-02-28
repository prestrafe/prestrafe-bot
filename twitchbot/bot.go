package twitchbot

import (
	"github.com/gempir/go-twitch-irc"

	"prestrafe-bot/config"
)

// This interface defines the public API of the chat bot. The API is pretty slim, as most of the work is done
// internally.
type ChatBot interface {
	// Joins a channel that is defined inside the passed channel configuration.
	Join(config *config.ChannelConfig) ChatBot
	// Starts up the chat bot in the thread that is calling this method. It blocks until an error occurs or the bot is
	// stopped.
	Start() error
	// Stops the bot and disconnects it from the Twitch API.
	Stop() error
}

type chatBot struct {
	channels []botChannel
	client   *twitch.Client
}

// Creates a new chat bot instance.
func NewChatBot(config *config.TwitchConfig) ChatBot {
	return &chatBot{
		channels: nil,
		client:   twitch.NewClient(config.UserName, config.AccessToken),
	}
}

func (c chatBot) Join(config *config.ChannelConfig) ChatBot {
	channel := newChannel(c.client, config)

	c.client.Join(channel.name)
	c.channels = append(c.channels, *channel)

	return c
}

func (c chatBot) Start() error {
	return c.client.Connect()
}

func (c chatBot) Stop() error {
	return c.client.Disconnect()
}
