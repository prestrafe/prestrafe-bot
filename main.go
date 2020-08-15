package main

import (
	"github.com/kelseyhightower/envconfig"

	"gitlab.com/prestrafe/prestrafe-bot/config"
	"gitlab.com/prestrafe/prestrafe-bot/globalapi"
	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
	"gitlab.com/prestrafe/prestrafe-bot/twitchbot"
)

type BotConfig struct {
	GlobalApiToken string `required:"true"`
	GsiAddr        string `required:"true"`
	GsiPort        int    `required:"true"`
	TwitchUsername string `required:"true"`
	TwitchApiToken string `required:"true"`
}

func main() {
	botConfig := new(BotConfig)
	envconfig.MustProcess("bot", botConfig)

	channelsConfig, configErr := config.ReadConfig("config.yml")
	if configErr != nil {
		panic("Could not read config file: " + configErr.Error())
	}

	bot := twitchbot.NewChatBot(botConfig.TwitchUsername, botConfig.TwitchApiToken)
	for _, channelConfig := range channelsConfig.Channels {
		bot.Join(channelConfig.Name, createCommands(botConfig, &channelConfig))
	}

	if twitchErr := bot.Start(); twitchErr != nil {
		panic("Twitch chat error: " + twitchErr.Error())
	}
}

func createCommands(botConfig *BotConfig, channelConfig *config.ChannelConfig) []twitchbot.ChatCommand {
	apiClient := globalapi.NewClient(botConfig.GlobalApiToken)
	gsiClient := gsiclient.New(botConfig.GsiAddr, botConfig.GsiPort, channelConfig.GsiToken)

	commands := []twitchbot.ChatCommand{
		// Troll commands
		twitchbot.NewGlobalCheckCommand().
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("globalcheck")).
			Build(),

		// Map information commands
		twitchbot.NewMapCommand(gsiClient, apiClient).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("map")).
			Build(),
		twitchbot.NewTierCommand(gsiClient, apiClient).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("tier")).
			Build(),

		// Player commands
		twitchbot.NewModeCommand(gsiClient).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("mode")).
			Build(),
		twitchbot.NewRankCommand(gsiClient, apiClient).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("rank")).
			Build(),
		twitchbot.NewStatsCommand(gsiClient).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("stats")).
			Build(),

		// Record time commands
		twitchbot.NewWRCommand(gsiClient, apiClient).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("wr")).
			Build(),
		twitchbot.NewBWRCommand(gsiClient, apiClient).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("wr")).
			Build(),
		twitchbot.NewPBCommand(gsiClient, apiClient).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("pb")).
			Build(),
		twitchbot.NewBPBCommand(gsiClient, apiClient).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("pb")).
			Build(),

		// Jump Stat commands
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "bh", "bhop", "Bunnyhop", 400, channelConfig.JumpsOnlyWithoutBinds).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("jumpstat")).
			Build(),
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "dh", "drophop", "Drop Bunnyhop", 400, channelConfig.JumpsOnlyWithoutBinds).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("jumpstat")).
			Build(),
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "laj", "ladderjump", "Ladder Jump", 400, channelConfig.JumpsOnlyWithoutBinds).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("jumpstat")).
			Build(),
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "lj", "longjump", "Long Jump", 300, channelConfig.JumpsOnlyWithoutBinds).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("jumpstat")).
			Build(),
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "mbh", "multibhop", "Multi Bunnyhop", 400, channelConfig.JumpsOnlyWithoutBinds).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("jumpstat")).
			Build(),
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "wj", "weirdjump", "Weird Jump", 400, channelConfig.JumpsOnlyWithoutBinds).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("jumpstat")).
			Build(),
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "lbh", "lbh", "Lowpre Bunnyhop", 400, channelConfig.JumpsOnlyWithoutBinds).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("jumpstat")).
			Build(),
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "lwj", "lwj", "Lowpre Weird Jump", 400, channelConfig.JumpsOnlyWithoutBinds).
			WithConfig(channelConfig.GetCommandConfig("*")).
			WithConfig(channelConfig.GetCommandConfig("jumpstat")).
			Build(),
	}

	commands = append(commands, twitchbot.NewHelpCommand(commands).
		WithConfig(channelConfig.GetCommandConfig("*")).
		WithConfig(channelConfig.GetCommandConfig("help")).
		Build(),
	)

	return commands
}
