package twitchbot

import (
	"testing"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/stretchr/testify/assert"
)

func TestCommandBuilderDefaults(t *testing.T) {
	command := NewChatCommandBuilder("name").build()

	if assert.NotNil(t, command) {
		assert.Equal(t, "name", command.name)
		assert.Empty(t, command.aliases)
		assert.Empty(t, command.parameters)
		assert.NotNil(t, command.handler)
	}
}

func TestCommandBuilder(t *testing.T) {
	command := NewChatCommandBuilder("name").
		WithAlias("alias1", "alias2").
		WithParameter("param", false, ".*").
		WithHandler(func(ctx CommandContext) (msg string, err error) {
			return "", nil
		}).
		build()

	if assert.NotNil(t, command) {
		assert.Equal(t, "name", command.name)
		assert.Equal(t, []string{"alias1", "alias2"}, command.aliases)

		if assert.Len(t, command.parameters, 1) {
			parameter := command.parameters[0]

			assert.Equal(t, "param", parameter.name)
			assert.False(t, parameter.required)
			assert.Equal(t, ".*", parameter.pattern)

			assert.Equal(t, "<param?>", parameter.String())
			assert.Equal(t, "(?P<param>.*)?", parameter.getPattern())
		}

		assert.NotNil(t, command.handler)
	}
}

func TestCommandBuilderWithConfig(t *testing.T) {
	command := NewChatCommandBuilder("name").build()

	if assert.NotNil(t, command) {
		assert.Equal(t, "name", command.name)
	}
}

func TestMatchesPrefix(t *testing.T) {
	command := NewChatCommandBuilder("name").
		WithAlias("alias1", "alias2").
		WithParameter("param", false, ".*").
		WithHandler(func(ctx CommandContext) (msg string, err error) {
			return "", nil
		}).
		build()

	assert.True(t, command.matchesPrefix("!name"))
	assert.True(t, command.matchesPrefix("!alias1"))
	assert.True(t, command.matchesPrefix("!alias2"))
}

func TestIsOnCoolDown(t *testing.T) {
	command := NewChatCommandBuilder("name").build()

	assert.False(t, command.isOnCoolDown())
	command.lastExecution = time.Now()
	assert.True(t, command.isOnCoolDown())
	command.lastExecution = time.Now().Add(-coolDown)
	assert.False(t, command.isOnCoolDown())
}

func TestParseParameters(t *testing.T) {
	noParamsCommand := NewChatCommandBuilder("name").build()
	assert.Empty(t, noParamsCommand.parseParameters("!name"))
	assert.Empty(t, noParamsCommand.parseParameters("!name ignored"))

	paramsCommand := NewChatCommandBuilder("name").WithParameter("param", false, "[0-9]+").build()
	assert.Empty(t, paramsCommand.parseParameters("!name"))
	assert.Empty(t, paramsCommand.parseParameters("!name ignored"))
	assert.NotEmpty(t, paramsCommand.parseParameters("!name 1337"))
}

func TestStringer(t *testing.T) {
	command := NewChatCommandBuilder("name").
		WithAlias("ignored").
		WithParameter("param1", true, ".*").
		WithParameter("param2", false, ".*").
		build()
	assert.Equal(t, "!name <param1> <param2?>", command.String())
}

func TestTryHandle(t *testing.T) {
	enabledCommand := NewChatCommandBuilder("name").
		WithParameter("param", true, "[0-9]+").
		build()
	abyss := func(format string, a ...interface{}) {}

	assert.False(t, enabledCommand.TryHandle("chan", &twitch.PrivateMessage{Message: ""}, abyss))
	assert.False(t, enabledCommand.TryHandle("chan", &twitch.PrivateMessage{Message: "!other"}, abyss))
	assert.False(t, enabledCommand.TryHandle("chan", &twitch.PrivateMessage{Message: "!other 42"}, abyss))

	assert.True(t, enabledCommand.TryHandle("chan", &twitch.PrivateMessage{Message: "!name"}, abyss))
	assert.True(t, enabledCommand.TryHandle("chan", &twitch.PrivateMessage{Message: "!name 42"}, abyss))
}

func TestCommandDetectionRegression(t *testing.T) {
	mapCommand := NewChatCommandBuilder("map").WithParameter("map", false, "[A-Za-z0-9_]+").build()
	pbCommand := NewChatCommandBuilder("pb").WithAlias("pr").WithParameter("map", false, "[A-Za-z0-9_]+").build()
	tierCommand := NewChatCommandBuilder("tier").WithParameter("map", false, "[A-Za-z0-9_]+").build()
	abyss := func(format string, a ...interface{}) {}

	assert.False(t, mapCommand.TryHandle("chan", &twitch.PrivateMessage{Message: "!mapcomp"}, abyss))
	assert.False(t, pbCommand.TryHandle("chan", &twitch.PrivateMessage{Message: "!prestrafebot"}, abyss))
	assert.False(t, tierCommand.TryHandle("chan", &twitch.PrivateMessage{Message: "!tiers"}, abyss))
}

func TestCommandDetectionWithSecondOptionalParameter(t *testing.T) {
	wrCommand := NewChatCommandBuilder("wr").
		WithAlias("gr", "gwr", "top").
		WithParameter("map", false, "(kz|kzpro|skz|vnl|xc)_[A-Za-z0-9_]+").
		WithParameter("mode", false, "(kzt|skz|vnl)").
		build()
	abyss := func(format string, a ...interface{}) {}

	assert.True(t, wrCommand.TryHandle("chan", &twitch.PrivateMessage{Message: "!wr skz"}, abyss))
}
