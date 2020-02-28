package old

import (
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc"

	"prestrafe-bot/config"
)

type Bot struct {
	config   *config.BotConfig
	client   *twitch.Client
	commands map[string]*Command
}

func New(botConfig *config.BotConfig) *Bot {
	client := twitch.NewClient(botConfig.Twitch.BotName, botConfig.Twitch.AccessToken)
	client.Join(botConfig.Twitch.ChannelName)

	return &Bot{
		config:   botConfig,
		client:   client,
		commands: map[string]*Command{},
	}
}

func (bot *Bot) AddCommand(name, configKey string, parameters int, handler CommandHandler) *Bot {
	bot.commands[name] = createCommand(bot.config.GetCommandConfig(configKey), parameters, handler)
	return bot
}

func (bot *Bot) Start() error {
	bot.client.OnNewMessage(bot.handleMessage)
	return bot.client.Connect()
}

func createCommand(config *config.ChatCommandConfig, parameters int, handler CommandHandler) *Command {
	return &Command{
		Enabled:    *config.Enabled,
		SubOnly:    *config.SubOnly,
		CoolDown:   time.Duration(*config.CoolDown) * time.Second,
		Parameters: parameters,
		Handler:    handler,
	}
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

type CommandHandler func(twitch.User, []string) string

type Command struct {
	Enabled  bool
	SubOnly  bool
	CoolDown time.Duration

	Parameters int
	Handler    CommandHandler

	lastExecution time.Time
}

func (cmd *Command) CanExecute(user twitch.User, arguments []string) bool {
	_, sub := user.Badges["subscriber"]
	_, mod := user.Badges["moderator"]
	_, broadcaster := user.Badges["broadcaster"]

	hasPermission := sub || mod || broadcaster || !cmd.SubOnly
	notInTimeout := cmd.lastExecution.Before(time.Now().Add(-cmd.CoolDown))

	return cmd.Enabled &&
		hasPermission &&
		notInTimeout &&
		len(arguments) >= cmd.Parameters
}

func (cmd *Command) Execute(user twitch.User, parameters []string) string {
	cmd.lastExecution = time.Now()
	return cmd.Handler(user, parameters)
}
