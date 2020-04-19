package twitchbot

import (
	"testing"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/stretchr/testify/assert"

	"gitlab.com/prestrafe/prestrafe-bot/config"
)

func TestCommandBuilderDefaults(t *testing.T) {
	command := NewChatCommandBuilder("name").build()

	if assert.NotNil(t, command) {
		assert.Equal(t, "name", command.name)
		assert.Empty(t, command.aliases)
		assert.Empty(t, command.parameters)
		assert.True(t, command.enabled)
		assert.False(t, command.subOnly)
		assert.Equal(t, 15*time.Second, command.coolDown)
		assert.NotNil(t, command.handler)
	}
}

func TestCommandBuilder(t *testing.T) {
	command := NewChatCommandBuilder("name").
		WithAlias("alias1", "alias2").
		WithParameter("param", false, ".*").
		WithEnabled(false).
		WithSubOnly(true).
		WithCoolDown(10 * time.Second).
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

		assert.False(t, command.enabled)
		assert.True(t, command.subOnly)
		assert.Equal(t, 10*time.Second, command.coolDown)
		assert.NotNil(t, command.handler)
	}
}

func TestCommandBuilderWithConfig(t *testing.T) {
	command := NewChatCommandBuilder("name").
		WithConfig(testCommandConfig(false, true, 10)).
		build()

	if assert.NotNil(t, command) {
		assert.Equal(t, "name", command.name)
		assert.False(t, command.enabled)
		assert.True(t, command.subOnly)
		assert.Equal(t, 10*time.Second, command.coolDown)
	}
}

func TestCommandBuilderWithConfigOverwrite(t *testing.T) {
	command := NewChatCommandBuilder("name").
		WithConfig(testCommandConfig(false, true, 10)).
		WithConfig(testCommandConfig(false, false, 20)).
		build()

	if assert.NotNil(t, command) {
		assert.Equal(t, "name", command.name)
		assert.False(t, command.enabled)
		assert.False(t, command.subOnly)
		assert.Equal(t, 20*time.Second, command.coolDown)
	}
}

func TestMatchesPrefix(t *testing.T) {
	command := NewChatCommandBuilder("name").
		WithAlias("alias1", "alias2").
		WithParameter("param", false, ".*").
		WithSubOnly(true).
		WithCoolDown(10 * time.Second).
		WithHandler(func(ctx CommandContext) (msg string, err error) {
			return "", nil
		}).
		build()

	assert.True(t, command.matchesPrefix("!name"))
	assert.True(t, command.matchesPrefix("!alias1"))
	assert.True(t, command.matchesPrefix("!alias2"))
}

func TestCanExecute(t *testing.T) {
	publicCommand := NewChatCommandBuilder("name").WithSubOnly(false).build()
	subOnlyCommand := NewChatCommandBuilder("name").WithSubOnly(true).build()

	normalUser := &twitch.User{Badges: map[string]int{}}
	subUser := &twitch.User{Badges: map[string]int{"subscriber": 1}}
	modUser := &twitch.User{Badges: map[string]int{"moderator": 1}}
	broadcasterUser := &twitch.User{Badges: map[string]int{"broadcaster": 1}}

	assert.True(t, publicCommand.canExecute(normalUser))
	assert.True(t, publicCommand.canExecute(subUser))
	assert.True(t, publicCommand.canExecute(modUser))
	assert.True(t, publicCommand.canExecute(broadcasterUser))
	assert.False(t, subOnlyCommand.canExecute(normalUser))
	assert.True(t, subOnlyCommand.canExecute(subUser))
	assert.True(t, subOnlyCommand.canExecute(modUser))
	assert.True(t, subOnlyCommand.canExecute(broadcasterUser))
}

func TestIsOnCoolDown(t *testing.T) {
	command := NewChatCommandBuilder("name").WithCoolDown(1 * time.Second).build()

	assert.False(t, command.isOnCoolDown())
	command.lastExecution = time.Now()
	assert.True(t, command.isOnCoolDown())
	command.lastExecution = time.Now().Add(-command.coolDown)
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
		WithSubOnly(true).
		build()
	normalUser := &twitch.User{Badges: map[string]int{}}
	abyss := func(format string, a ...interface{}) {}

	assert.False(t, enabledCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: ""}, abyss))
	assert.False(t, enabledCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: "!other"}, abyss))
	assert.False(t, enabledCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: "!other 42"}, abyss))

	assert.True(t, enabledCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: "!name"}, abyss))
	assert.True(t, enabledCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: "!name 42"}, abyss))

	disabledCommand := NewChatCommandBuilder("name").
		WithParameter("param", true, "[0-9]+").
		WithEnabled(false).
		build()

	assert.False(t, disabledCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: ""}, abyss))
	assert.False(t, disabledCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: "!other"}, abyss))
	assert.False(t, disabledCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: "!other 42"}, abyss))

	assert.False(t, disabledCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: "!name"}, abyss))
	assert.False(t, disabledCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: "!name 42"}, abyss))
}

func TestCommandDetectionRegression(t *testing.T) {
	mapCommand := NewChatCommandBuilder("map").WithParameter("map", false, "[A-Za-z0-9_]+").build()
	pbCommand := NewChatCommandBuilder("pb").WithAlias("pr").WithParameter("map", false, "[A-Za-z0-9_]+").build()
	tierCommand := NewChatCommandBuilder("tier").WithParameter("map", false, "[A-Za-z0-9_]+").build()
	normalUser := &twitch.User{Badges: map[string]int{}}
	abyss := func(format string, a ...interface{}) {}

	assert.False(t, mapCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: "!mapcomp"}, abyss))
	assert.False(t, pbCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: "!prestrafebot"}, abyss))
	assert.False(t, tierCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: "!tiers"}, abyss))
}

func TestCommandDetectionWithSecondOptionalParameter(t *testing.T) {
	wrCommand := NewChatCommandBuilder("wr").
		WithAlias("gr", "gwr", "top").
		WithParameter("map", false, "(kz|kzpro|skz|vnl|xc)_[A-Za-z0-9_]+").
		WithParameter("mode", false, "(kzt|skz|vnl)").
		build()

	normalUser := &twitch.User{Badges: map[string]int{}}
	abyss := func(format string, a ...interface{}) {}

	assert.True(t, wrCommand.TryHandle("chan", normalUser, &twitch.PrivateMessage{Message: "!wr skz"}, abyss))
}

func testCommandConfig(enabled bool, subOnly bool, coolDown int) *config.ChatCommandConfig {
	return &config.ChatCommandConfig{
		Name:     "",
		Enabled:  &enabled,
		SubOnly:  &subOnly,
		CoolDown: &coolDown,
	}
}
