package twitchbot

import (
	"fmt"

	"github.com/gempir/go-twitch-irc"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func HandlePBCommand(user twitch.User, strings []string) string {
	gameState, err := gsi.GetGameState()
	if err != nil {
		return "Could not retrieve current game state"
	}

	nub, pro, err := globalapi.GetPersonalBests(gameState.Map.Name, gameState.Player.TimerMode(), gameState.Player.SteamId)

	message := fmt.Sprintf("%s on %s [%s]: ", gameState.Player.Name, gameState.Map.Name, gameState.Player.Clan)
	if nub != nil {
		message += fmt.Sprintf("NUB: %s (%d TP)", nub.FormattedTime(), nub.Teleports)
	} else {
		message += fmt.Sprintf("NUB: None")
	}

	message += ", "

	if pro != nil {
		message += fmt.Sprintf("PRO: %s", pro.FormattedTime())
	} else {
		message += fmt.Sprintf("PRO: None")
	}

	return message
}
