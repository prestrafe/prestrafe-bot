package main

import (
	"gitlab.com/prestrafe/prestrafe-bot/config"
	"gitlab.com/prestrafe/prestrafe-bot/gsi"
	"gitlab.com/prestrafe/prestrafe-bot/twitchbot"
)

func main() {
	botConfig, configErr := config.ReadConfig("config.yaml")
	if configErr != nil {
		panic("Could not read config file: " + configErr.Error())
	}

	go runGameStateIntegration(botConfig)
	go runTwitch(botConfig)

	/*
		callback := func(channel string, user twitch.User, message twitch.Message) {
			fmt.Println(user, message)
		}

		client := twitch.NewClient(botConfig.Twitch.BotName, botConfig.Twitch.AccessToken)
		client.Join("nykan")
		client.OnNewUsernoticeMessage(callback)
		client.OnNewNoticeMessage(callback)
		client.Connect()
	*/

	select {}
}

func runGameStateIntegration(botConfig *config.BotConfig) {
	s := gsi.NewServer(&botConfig.Gsi)
	if err := s.Start(); err != nil {
		panic(err)
	}
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
