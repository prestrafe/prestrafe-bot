package twitchbot

import (
	"errors"
	"fmt"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func NewTierCommand(gsiClient gsi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("tier").
		WithAlias("difficulty").
		WithParameter("map", false, "[A-Za-z0-9_]").
		WithHandler(createTierHandler(gsiClient))
}

func createTierHandler(gsiClient gsi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		mapName, hasMapName := ctx.Parameter("map")

		if !hasMapName {
			gameState, gsiError := gsiClient.GetGameState()
			if gsiError != nil {
				return "", errors.New("could not retrieve KZ game play")
			}

			mapName = gameState.Map.Name
		}

		globalMap, apiError := globalapi.GetMapByName(mapName)
		if apiError != nil {
			return fmt.Sprintf("%s - Tier Unknown (Not global)", mapName), nil
		} else {
			return fmt.Sprintf("%s - Tier %d", mapName, globalMap.Difficulty), nil
		}
	}
}
