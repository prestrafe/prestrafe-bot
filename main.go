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
	s := gsi.CreateServer(botConfig.Server.VerificationToken, time.Duration(botConfig.Server.TTL)*time.Second)
	if err := s.ListenAndServer(); err != nil {
		panic(err)
	}
}

func runTwitch(botConfig *config.BotConfig) {
	bot := twitchbot.New(botConfig)

	bot.AddCommand("map", "map", 0, twitchbot.HandleMapCommand)
	bot.AddCommand("wr", "wr", 0, twitchbot.HandleWRCommand)
	bot.AddCommand("bwr", "wr", 0, twitchbot.HandleBonusWRCommand)
	bot.AddCommand("pb", "pb", 0, twitchbot.HandlePBCommand)
	bot.AddCommand("bpb", "pb", 0, twitchbot.HandleBonusPBCommand)

	bot.AddCommand("bh", "js", 0, twitchbot.CreateJSHandler("bhop", "Bunnyhop"))
	bot.AddCommand("dh", "js", 0, twitchbot.CreateJSHandler("drophop", "Drop Bunnyhop"))
	bot.AddCommand("laj", "js", 0, twitchbot.CreateJSHandler("ladderjump", "Ladder Jump"))
	bot.AddCommand("lj", "js", 0, twitchbot.CreateJSHandler("longjump", "Long Jump"))
	bot.AddCommand("mbh", "js", 0, twitchbot.CreateJSHandler("multibhop", "Multi Bunnyhop"))
	bot.AddCommand("wj", "js", 0, twitchbot.CreateJSHandler("weirdjump", "Weird Jump"))

	bot.AddCommand("globalcheck", "globalcheck", 0, twitchbot.HandleGlobalCheckCommand)

	bot.AddCommand("kz", "help", 0, bot.CreateHelpCommand())

	if twitchErr := bot.Start(); twitchErr != nil {
		panic("Twitch chat error: " + twitchErr.Error())
	}
}
