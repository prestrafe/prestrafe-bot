package twitchbot

import (
	"strings"

	"github.com/gempir/go-twitch-irc"
)

type Bot struct {
	channel  string
	client   *twitch.Client
	commands map[string]*Command
}

func New(channel, name, token string) *Bot {
	client := twitch.NewClient(name, token)
	client.Join(channel)

	return &Bot{
		channel:  channel,
		client:   client,
		commands: map[string]*Command{},
	}
}

func (bot *Bot) AddCommand(name string, command *Command) *Bot {
	bot.commands[name] = command
	return bot
}

func (bot *Bot) Start() error {
	bot.client.OnNewMessage(bot.handleMessage)
	return bot.client.Connect()
}

func (bot *Bot) handleMessage(channel string, user twitch.User, message twitch.Message) {
	commandName, parameters := parseCommand(message.Text)
	if command, contains := bot.commands[commandName]; contains && command.CanExecute(user, parameters) {
		bot.client.Say(channel, command.Execute(user, parameters))
	}
}

func parseCommand(message string) (name string, parameters []string) {
	if strings.HasPrefix(message, "!") {
		parts := strings.Split(message[1:], " ")
		name = parts[0]
		parameters = parts[1:]
	}

	return
}
