package main

import (
	"gitlab.com/prestrafe/prestrafe-bot/config"
	"gitlab.com/prestrafe/prestrafe-bot/twitchbot"
)

func main() {
	botConfig, configErr := config.ReadConfig("config.yml")
	if configErr != nil {
		panic("Could not read config file: " + configErr.Error())
	}

	runTwitch(botConfig)
}

func runTwitch(botConfig *config.BotConfig) {
	bot := twitchbot.NewChatBot(&botConfig.Twitch, &botConfig.Gsi)
	for _, channelConfig := range botConfig.Channels {
		bot.Join(&channelConfig)
	}

	if twitchErr := bot.Start(); twitchErr != nil {
		panic("Twitch chat error: " + twitchErr.Error())
	}
}
