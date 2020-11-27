package twitchbot

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
)

const (
	commandPrefix = "!"
	coolDown      = 15 * time.Second
)

// Defines function signatures that can handle chat commands.
// They get passed a the parameters that were passed to the command by the author.
// Handlers can respond with either a message or an error. Both will get printed to chat, but an error does not cause
// the command to be cooled down as much.
type ChatCommandHandler func(ctx CommandContext) (msg string, err error)

// Defines a function signature that can be called to send chat messages to an unknown source. This way an abstraction
// is provided between the command handling and the actual receiver of command results.
type ChatMessageSink func(format string, a ...interface{})

// Defines an interface that exposes the external API of a chat command. A chat command encapsulates a lot of common
// logic for all command handlers and takes care of permissions, cool downs, error handling and message parsing.
type ChatCommand interface {
	// Returns the name of the command.
	Name() string
	// Tries to handle a chat message from a chat user. If the command is able to identify that the message is handled
	// by it, this method returns true, regardless if the command later produces an error. It returns false only if the
	// command does not claim responsibility for the message.
	TryHandle(channel string, message *twitch.PrivateMessage, messageSink ChatMessageSink) bool
	// Returns a string representation of the command. This is usually the signature of the command.
	String() string
}

type chatCommand struct {
	name          string
	aliases       []string
	parameters    []chatCommandParameter
	handler       ChatCommandHandler
	pattern       *regexp.Regexp
	lastExecution time.Time
}

func (c *chatCommand) Name() string {
	return c.name
}

func (c *chatCommand) TryHandle(channel string, message *twitch.PrivateMessage, messageSink ChatMessageSink) bool {
	if !c.matchesPrefix(message.Message) {
		return false
	}

	if !c.pattern.MatchString(message.Message) {
		messageSink("Your message did not matched the usage of the command: %s", c)
		return true
	}

	if c.isOnCoolDown() {
		messageSink("This command has a cool down of %.0f seconds. Please try again later.", coolDown.Seconds())
		return true
	}

	output, err := c.handler(&commandContext{channel, c.parseParameters(message.Message)})
	if err != nil {
		messageSink(err.Error())
	} else {
		messageSink(output)
	}

	c.lastExecution = time.Now()
	return true
}

func (c *chatCommand) String() string {
	var parameters []string
	for _, p := range c.parameters {
		parameters = append(parameters, p.String())
	}

	parameterSignature := ""
	if len(parameters) > 0 {
		parameterSignature += " " + strings.Join(parameters, " ")
	}

	return fmt.Sprintf("%s%s%s", commandPrefix, c.name, parameterSignature)
}

func (c *chatCommand) matchesPrefix(message string) bool {
	namePattern := strings.Join(append(c.aliases, c.name), "|")
	commandPattern := fmt.Sprintf("^%s(%s)((\\s.+)|$)", commandPrefix, namePattern)

	matchString, _ := regexp.MatchString(commandPattern, message)
	return matchString
}

func (c *chatCommand) isOnCoolDown() bool {
	return c.lastExecution.After(time.Now().Add(-coolDown))
}

func (c *chatCommand) parseParameters(message string) map[string]string {
	match := c.pattern.FindStringSubmatch(message)
	groups := make(map[string]string)

	for i, name := range c.pattern.SubexpNames() {
		if i > 0 && i <= len(match) && len(name) > 0 && len(match[i]) > 0 {
			groups[name] = match[i]
		}
	}

	return groups
}
