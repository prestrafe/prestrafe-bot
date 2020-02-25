package twitchbot

import (
	"fmt"
	"strconv"

	"github.com/gempir/go-twitch-irc"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func HandleBonusPBCommand(user twitch.User, parameters []string) string {
	gameState, err := gsi.GetGameState()
	if err != nil {
		return "Could not retrieve KZ gameplay"
	}

	stage := 1
	if parsedStage, stageErr := strconv.Atoi(parameters[0]); stageErr == nil && parsedStage > 0 {
		stage = parsedStage
	}

	nub, pro, err := globalapi.GetPersonalRecord(gameState.Map.Name, gameState.Player.TimerMode(), stage, gameState.Player.SteamId)

	message := fmt.Sprintf("PB of %s on %s Bonus %d [%s]: ", gameState.Player.Name, gameState.Map.Name, stage, gameState.Player.Clan)
	if nub != nil && err == nil {
		message += fmt.Sprintf("NUB: %s (%d TP)", nub.FormattedTime(), nub.Teleports)
	} else {
		message += fmt.Sprintf("NUB: None")
	}

	message += ", "

	if pro != nil && err == nil {
		message += fmt.Sprintf("PRO: %s", pro.FormattedTime())
	} else {
		message += fmt.Sprintf("PRO: None")
	}

	return message
}
