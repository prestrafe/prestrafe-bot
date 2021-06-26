package twitchbot

import (
	"errors"
	"fmt"

	"gitlab.com/prestrafe/prestrafe-bot/globalapi"
	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
)

func NewMapCommand(gsiClient gsiclient.Client, apiClient globalapi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("map").
		WithParameter("map", false, mapRegexPattern).
		WithHandler(createMapHandler(gsiClient, apiClient))
}

func createMapHandler(gsiClient gsiclient.Client, apiClient globalapi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		var modeName string
		mapName, hasMapName := ctx.Parameter("map")

		gameState, gsiError := gsiClient.GetGameState()
		if !hasMapName && (gsiError != nil || !gsiclient.IsKZGameState(gameState)) {
			return "", errors.New("could not retrieve KZ game play")
		} else {
			modeName = gsiclient.TimerMode(gameState.Player)
		}

		if !hasMapName {
			mapName = gsiclient.GetMapName(gameState.Map)
		}

		globalMap, apiError := (&globalapi.MapServiceClient{Client: apiClient}).GetMapByName(mapName)
		if globalMap == nil || apiError != nil {
			return fmt.Sprintf("Map: %s (Not global)", mapName), nil
		} else {
			return fmt.Sprintf("Map: %s (T%d) - https://gokzstats.com/?map=%s&mode=%s", mapName, globalMap.Difficulty, mapName, modeName), nil
		}
	}
}
