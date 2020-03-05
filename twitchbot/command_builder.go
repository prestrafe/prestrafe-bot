package twitchbot

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gitlab.com/prestrafe/prestrafe-bot/config"
)

// The command builder is used to create new command definitions in a declarative, readable and type-safe way.
// It also provides default values for all builder methods, so only passing a name when eventually building a command
// instance is actually required.
type ChatCommandBuilder interface {
	// Adds zero or more aliases to the command that is being build. An alias defines additional names that can be used
	// to trigger the build command as if it was trigger by its real name. By default a command will have no aliases
	// defined.
	WithAlias(alias ...string) ChatCommandBuilder
	// Adds a new parameter to the command that is being build. A parameter allows additional information being passed
	// to the command handler. A parameter is defined by three properties:
	// 	- A name, that is used to identify it in help texts and by command handlers.
	// 	- The fact that it is required or not.
	// 	- A regular expression pattern that needs to be matched by values of the parameter.
	// By default no parameters will be defined.
	WithParameter(name string, required bool, pattern string) ChatCommandBuilder
	// Applies a chat command configuration element onto the builder. Calling this method may overwrite any previously
	// set configuration values of the receiving builder, but future calls may also overwrite parts of the passed
	// configuration. It is also possible to call this method more then once to layer configurations on top of one
	// another.
	WithConfig(config *config.ChatCommandConfig) ChatCommandBuilder
	// Sets the command that is being build to be enabled or disabled. By default the command will be enabled.
	WithEnabled(enabled bool) ChatCommandBuilder
	// Sets the command that is being build to be only available to subscribers, moderators and broadcasters. By default
	// the command will be available to everyone in chat.
	WithSubOnly(subOnly bool) ChatCommandBuilder
	// Sets the cool down duration for the command that is being build. The cool down ensures that at least the duration
	// of the cool down has passed between two durations of the command.
	WithCoolDown(coolDown time.Duration) ChatCommandBuilder
	// Sets the handler of the command that is being build. The handler is invoked every time the command is triggered
	// by a chat user. By default the handler will display a message that the command is not yet implemented correctly.
	WithHandler(handler ChatCommandHandler) ChatCommandBuilder
	// Will build a new command instance, initialized with the values passed to the builder instance.
	Build() ChatCommand
	build() *chatCommand
}

// Creates a new command builder. The builder will already be initialized with the command name, as setting a name is
// required.
func NewChatCommandBuilder(name string) ChatCommandBuilder {
	return &chatCommand{
		name,
		[]string{},
		[]chatCommandParameter{},
		true,
		false,
		15 * time.Second,
		func(ctx CommandContext) (msg string, err error) {
			return "", errors.New("this command is not yet implemented")
		},
		nil,
		time.Unix(0, 0),
	}
}

func (c *chatCommand) WithAlias(alias ...string) ChatCommandBuilder {
	c.aliases = append(c.aliases, alias...)
	return c
}

func (c *chatCommand) WithParameter(name string, required bool, pattern string) ChatCommandBuilder {
	c.parameters = append(c.parameters, chatCommandParameter{
		name:     name,
		required: required,
		pattern:  pattern,
	})
	return c
}

func (c *chatCommand) WithConfig(config *config.ChatCommandConfig) ChatCommandBuilder {
	if config.Enabled != nil {
		c.enabled = *config.Enabled
	}
	if config.SubOnly != nil {
		c.subOnly = *config.SubOnly
	}
	if config.CoolDown != nil {
		c.coolDown = time.Duration(*config.CoolDown) * time.Second
	}

	return c
}

func (c *chatCommand) WithEnabled(enabled bool) ChatCommandBuilder {
	c.enabled = enabled
	return c
}

func (c *chatCommand) WithSubOnly(subOnly bool) ChatCommandBuilder {
	c.subOnly = subOnly
	return c
}

func (c *chatCommand) WithCoolDown(coolDown time.Duration) ChatCommandBuilder {
	c.coolDown = coolDown
	return c
}

func (c *chatCommand) WithHandler(handler ChatCommandHandler) ChatCommandBuilder {
	c.handler = handler
	return c
}

func (c *chatCommand) Build() ChatCommand {
	return c.build()
}

func (c *chatCommand) build() *chatCommand {
	namePattern := strings.Join(append(c.aliases, c.name), "|")
	commandPattern := fmt.Sprintf("%s(%s)", commandPrefix, namePattern)

	for _, p := range c.parameters {
		commandPattern += fmt.Sprintf("(\\s+|$)%s", p.getPattern())
	}

	c.pattern = regexp.MustCompile(commandPattern)
	c.lastExecution = time.Unix(0, 0)
	return c
}
