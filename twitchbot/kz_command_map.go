package twitchbot

import (
	"errors"
	"fmt"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func NewMapCommand(gsiClient gsi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("map").
		WithParameter("map", false, "[A-Za-z0-9_]").
		WithHandler(createMapHandler(gsiClient))
}

func createMapHandler(gsiClient gsi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		mapName, hasMapName := ctx.Parameter("map")

		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil {
			return "", errors.New("could not retrieve KZ game play")
		}

		if !hasMapName {
			mapName = gameState.Map.Name
		}

		globalMap, apiError := globalapi.GetMapByName(mapName)
		if apiError != nil {
			return fmt.Sprintf("Map: %s (Not global)", mapName), nil
		} else {
			return fmt.Sprintf("Map: %s (T%d) - https://gokzstats.com/?map=%s&mode=%s", mapName, globalMap.Difficulty, mapName, gameState.Player.TimerMode()), nil
		}
	}
}
