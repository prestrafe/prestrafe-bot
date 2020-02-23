package twitchbot

import (
	"time"

	"github.com/gempir/go-twitch-irc"

	"prestrafe-bot/config"
)

type CommandHandler func(twitch.User, []string) string

type Command struct {
	Enabled  bool
	SubOnly  bool
	CoolDown time.Duration

	Parameters int
	Handler    CommandHandler

	lastExecution time.Time
}

func CreateCommand(config *config.CommandConfig, parameters int, handler CommandHandler) *Command {
	return &Command{
		Enabled:    *config.Enabled,
		SubOnly:    *config.SubOnly,
		CoolDown:   time.Duration(*config.CoolDown) * time.Second,
		Parameters: parameters,
		Handler:    handler,
	}
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
