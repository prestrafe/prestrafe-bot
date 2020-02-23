package twitchbot

import (
	"fmt"

	"github.com/gempir/go-twitch-irc"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func HandleWRCommand(user twitch.User, strings []string) string {
	gameState, err := gsi.GetGameState()
	if err != nil {
		return "Could not retrieve current game state"
	}

	nub, pro, err := globalapi.GetWorldRecord(gameState.Map.Name, gameState.Player.TimerMode())

	message := fmt.Sprintf("Global Records on %s [%s]: ", gameState.Map.Name, gameState.Player.Clan)
	if nub != nil {
		message += fmt.Sprintf("NUB: %s (%d TP) by %s", nub.FormattedTime(), nub.Teleports, nub.PlayerName)
	} else {
		message += fmt.Sprintf("NUB: None")
	}

	message += ", "

	if pro != nil {
		message += fmt.Sprintf("PRO: %s by %s", pro.FormattedTime(), pro.PlayerName)
	} else {
		message += fmt.Sprintf("PRO: None")
	}

	return message
}
