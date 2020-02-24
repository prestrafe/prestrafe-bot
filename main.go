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
	bot := twitchbot.New(botConfig)

	bot.AddCommand("map", "map", 0, twitchbot.HandleMapCommand)
	bot.AddCommand("wr", "wr", 0, twitchbot.HandleWRCommand)
	bot.AddCommand("pb", "pb", 0, twitchbot.HandlePBCommand)

	bot.AddCommand("bhpb", "js", 0, twitchbot.CreateJSHandler("bhop", "Bunnyhop"))
	bot.AddCommand("dhpb", "js", 0, twitchbot.CreateJSHandler("drophop", "Drop Bunnyhop"))
	bot.AddCommand("lajpb", "js", 0, twitchbot.CreateJSHandler("ladderjump", "Ladder Jump"))
	bot.AddCommand("ljpb", "js", 0, twitchbot.CreateJSHandler("longjump", "Long Jump"))
	bot.AddCommand("mbhpb", "js", 0, twitchbot.CreateJSHandler("multibhop", "Multi Bunnyhop"))
	bot.AddCommand("wjpb", "js", 0, twitchbot.CreateJSHandler("weirdjump", "Weird Jump"))

	if twitchErr := bot.Start(); twitchErr != nil {
		panic("Twitch chat error: " + twitchErr.Error())
	}
}
