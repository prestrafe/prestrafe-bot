package twitchbot

import (
	"strings"

	"github.com/gempir/go-twitch-irc"

	"gitlab.com/prestrafe/prestrafe-bot/config"
	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
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
	gsiConfig *config.GsiConfig
	channels  map[string]botChannel
	client    *twitch.Client
}

// Creates a new chat bot instance.
func NewChatBot(twitchConfig *config.TwitchConfig, gsiConfig *config.GsiConfig) ChatBot {
	return &chatBot{
		gsiConfig,
		make(map[string]botChannel),
		twitch.NewClient(twitchConfig.UserName, twitchConfig.AccessToken),
	}
}

func (c chatBot) Join(config *config.ChannelConfig) ChatBot {
	gsiClient := gsiclient.New(c.gsiConfig.Addr, c.gsiConfig.Port, config.GsiToken)

	channel := newChannel(c.client, gsiClient, config)
	channelName := strings.ToLower(channel.name)

	c.client.Join(channelName)
	c.channels[channelName] = *channel

	return c
}

func (c chatBot) Start() error {
	c.client.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		if botChannel, hasChannel := c.channels[channel]; hasChannel {
			botChannel.handle(&user, &message)
		}
	})
	return c.client.Connect()
}

func (c chatBot) Stop() error {
	return c.client.Disconnect()
}
