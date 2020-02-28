package twitchbot

import (
	"errors"
	"fmt"

	"prestrafe-bot/config"
	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func CreateWrCommand(config *config.ChatCommandConfig) ChatCommand {
	return NewChatCommandBuilder("wr").
		WithConfig(config).
		WithAlias("gr", "gwr", "top").
		WithParameter("map", false, "[A-Za-z0-9_]").
		WithHandler(handleWrCommand).
		Build()
}

func handleWrCommand(parameters map[string]string) (message string, err error) {
	mapName, hasMapName := parameters["map"]

	gameState, gsiError := gsi.GetGameState()
	if gsiError != nil {
		return "", errors.New("could not retrieve KZ game play")
	}

	if !hasMapName {
		mapName = gameState.Map.Name
	}

	nub, pro, apiError := globalapi.GetWorldRecord(mapName, gameState.Player.TimerMode(), 0)

	message = fmt.Sprintf("Global Records on %s [%s]: ", mapName, gameState.Player.Clan)
	if nub != nil && apiError == nil {
		message += fmt.Sprintf("NUB: %s (%d TP) by %s", nub.FormattedTime(), nub.Teleports, nub.PlayerName)
	} else {
		message += fmt.Sprintf("NUB: None")
	}

	message += ", "

	if pro != nil && apiError == nil {
		message += fmt.Sprintf("PRO: %s by %s", pro.FormattedTime(), pro.PlayerName)
	} else {
		message += fmt.Sprintf("PRO: None")
	}

	return
}
