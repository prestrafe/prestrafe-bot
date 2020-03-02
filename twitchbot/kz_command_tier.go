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
		WithParameter("map", false, "[A-Za-z0-9_]+").
		WithHandler(createTierHandler(gsiClient))
}

func createTierHandler(gsiClient gsi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		mapName, hasMapName := ctx.Parameter("map")

		if !hasMapName {
			gameState, gsiError := gsiClient.GetGameState()
			if gsiError != nil || !gameState.IsKZGameState() {
				return "", errors.New("could not retrieve KZ game play")
			}

			mapName = gameState.Map.GetMapName()
		}

		globalMap, apiError := globalapi.GetMapByName(mapName)
		if globalMap == nil || apiError != nil {
			return fmt.Sprintf("%s - Tier Unknown (Not global)", mapName), nil
		} else {
			return fmt.Sprintf("%s - Tier %d (%s)", mapName, globalMap.Difficulty, convertDifficultyToTier(globalMap.Difficulty)), nil
		}
	}
}

func convertDifficultyToTier(difficulty int) string {
	switch difficulty {
	case 1:
		return "Very Easy"
	case 2:
		return "Easy"
	case 3:
		return "Medium"
	case 4:
		return "Hard"
	case 5:
		return "Very Hard"
	case 6:
		return "Death"
	}

	return "Unknown"
}
