package main

import (
	"fmt"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gitlab.com/prestrafe/prestrafe-bot/config"
	"gitlab.com/prestrafe/prestrafe-bot/globalapi"
	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
	"gitlab.com/prestrafe/prestrafe-bot/smclient"
	"gitlab.com/prestrafe/prestrafe-bot/twitchbot"
)

type BotConfig struct {
	GlobalApiToken string `required:"true"`
	GsiAddr        string `required:"true"`
	GsiPort        int    `required:"true"`
	SmAddr         string `required:"true"`
	SmPort         int    `required:"true"`
	TwitchUsername string `required:"true"`
	TwitchApiToken string `required:"true"`
	ConfigDir      string `default:""`
	MetricPort     int    `default:"9080"`
}

func main() {
	botConfig := new(BotConfig)
	envconfig.MustProcess("bot", botConfig)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%d", botConfig.MetricPort), nil)
	}()

	channelsConfig, configErr := config.ReadConfig(fmt.Sprintf("%s/config.yml", botConfig.ConfigDir))
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
	smClient := smclient.New(botConfig.SmAddr, botConfig.SmPort, channelConfig.ServerToken)

	commands := []twitchbot.ChatCommand{
		// Globalcheck command
		twitchbot.NewGlobalCheckCommand(gsiClient, smClient).Build(),

		// Current run command
		twitchbot.NewRunCommand(gsiClient, smClient).Build(),

		// Server command
		twitchbot.NewServerCommand(gsiClient, smClient).Build(),

		// Map information commands
		twitchbot.NewMapCommand(gsiClient, apiClient).Build(),
		twitchbot.NewTierCommand(gsiClient, apiClient).Build(),

		// Player commands
		twitchbot.NewModeCommand(gsiClient).Build(),
		twitchbot.NewRankCommand(gsiClient, apiClient).Build(),
		twitchbot.NewStatsCommand(gsiClient).Build(),

		// Record time commands
		twitchbot.NewWRCommand(gsiClient, apiClient).Build(),
		twitchbot.NewBWRCommand(gsiClient, apiClient).Build(),
		twitchbot.NewPBCommand(gsiClient, apiClient).Build(),
		twitchbot.NewBPBCommand(gsiClient, apiClient).Build(),

		// Jump Stat commands
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "bh", "bhop", "Bunnyhop", 400).Build(),
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "dh", "drophop", "Drop Bunnyhop", 400).Build(),
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "laj", "ladderjump", "Ladder Jump", 400).Build(),
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "lj", "longjump", "Long Jump", 300).Build(),
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "mbh", "multibhop", "Multi Bunnyhop", 400).Build(),
		twitchbot.NewJumpStatCommand(gsiClient, apiClient, "wj", "weirdjump", "Weird Jump", 400).Build(),
	}

	commands = append(commands, twitchbot.NewHelpCommand(commands).Build())

	return commands
}
