package main

import (
	"time"

	"prestrafe-bot/config"
	"prestrafe-bot/gsi"
	"prestrafe-bot/twitchbot"
)

func main() {
	botConfig, configErr := config.ReadConfig("config.yaml")
	if configErr != nil {
		panic("Could not read config file: " + configErr.Error())
	}

	go runGameStateIntegration(botConfig)
	go runTwitch(botConfig)

	select {}
}

func runGameStateIntegration(botConfig *config.BotConfig) {
	s := gsi.CreateServer(time.Duration(botConfig.Server.TTL) * time.Second)
	if err := s.ListenAndServer(); err != nil {
		panic(err)
	}
}

func runTwitch(botConfig *config.BotConfig) {
	bot := twitchbot.New(botConfig.Twitch.ChannelName, botConfig.Twitch.BotName, botConfig.Twitch.AccessToken)

	mapConfig := botConfig.GetCommandConfig("map")
	bot.AddCommand("map", twitchbot.CreateCommand(mapConfig, 0, twitchbot.HandleMapCommand))

	wrConfig := botConfig.GetCommandConfig("wr")
	bot.AddCommand("wr", twitchbot.CreateCommand(wrConfig, 0, twitchbot.HandleWRCommand))

	pbConfig := botConfig.GetCommandConfig("pb")
	bot.AddCommand("pb", twitchbot.CreateCommand(pbConfig, 0, twitchbot.HandlePBCommand))

	jsCommandsConfig := botConfig.GetCommandConfig("js")
	bot.AddCommand("bhpb", twitchbot.CreateCommand(jsCommandsConfig, 0, twitchbot.CreateHandleJSCommand("bhop", "Bunnyhop")))
	bot.AddCommand("ljpb", twitchbot.CreateCommand(jsCommandsConfig, 0, twitchbot.CreateHandleJSCommand("longjump", "Long Jump")))

	if twitchErr := bot.Start(); twitchErr != nil {
		panic("Twitch chat error: " + twitchErr.Error())
	}
}
