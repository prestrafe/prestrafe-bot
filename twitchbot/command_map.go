package twitchbot

import (
	"fmt"

	"github.com/gempir/go-twitch-irc"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func HandleMapCommand(user twitch.User, parameters []string) string {
	gameState, err := gsi.GetGameState()
	if err != nil {
		return "Could not retrieve current game state"
	}

	globalMap, err := globalapi.GetMapByName(gameState.Map.Name)
	if err != nil {
		return fmt.Sprintf("Current map: %s (Not global)", gameState.Map.Name)
	} else {
		return fmt.Sprintf("Current map: %s (T%d) - https://gokzstats.com/?map=%s&mode=%s", gameState.Map.Name, globalMap.Difficulty, gameState.Map.Name, gameState.Player.TimerMode())
	}
}
