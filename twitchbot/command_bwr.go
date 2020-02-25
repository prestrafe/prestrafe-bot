package twitchbot

import (
	"fmt"
	"strconv"

	"github.com/gempir/go-twitch-irc"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func HandleBonusWRCommand(user twitch.User, parameters []string) string {
	gameState, err := gsi.GetGameState()
	if err != nil {
		return "Could not retrieve KZ gameplay"
	}

	stage := 1
	if parsedStage, stageErr := strconv.Atoi(parameters[0]); stageErr == nil && parsedStage > 0 {
		stage = parsedStage
	} else {
		return fmt.Sprintf("'%s' is not a valid bonus number.", parameters[0])
	}

	nub, pro, err := globalapi.GetWorldRecord(gameState.Map.Name, gameState.Player.TimerMode(), stage)

	message := fmt.Sprintf("Global Records on %s Bonus %d [%s]: ", gameState.Map.Name, stage, gameState.Player.Clan)
	if nub != nil && err == nil {
		message += fmt.Sprintf("NUB: %s (%d TP) by %s", nub.FormattedTime(), nub.Teleports, nub.PlayerName)
	} else {
		message += fmt.Sprintf("NUB: None")
	}

	message += ", "

	if pro != nil && err == nil {
		message += fmt.Sprintf("PRO: %s by %s", pro.FormattedTime(), pro.PlayerName)
	} else {
		message += fmt.Sprintf("PRO: None")
	}

	return message
}
