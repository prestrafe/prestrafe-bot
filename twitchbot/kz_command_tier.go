package twitchbot

import (
	"errors"
	"fmt"

	"gitlab.com/prestrafe/prestrafe-bot/globalapi"
	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
)

func NewTierCommand(gsiClient gsiclient.Client, apiClient globalapi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("tier").
		WithAlias("difficulty").
		WithParameter("map", false, "(kz|kzpro|skz|vnl|xc)_[A-Za-z0-9_]+").
		WithHandler(createTierHandler(gsiClient, apiClient))
}

func createTierHandler(gsiClient gsiclient.Client, apiClient globalapi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		mapName, hasMapName := ctx.Parameter("map")

		if !hasMapName {
			gameState, gsiError := gsiClient.GetGameState()
			if gsiError != nil || !gsiclient.IsKZGameState(gameState) {
				return "", errors.New("could not retrieve KZ game play")
			}

			mapName = gsiclient.GetMapName(gameState.Map)
		}

		globalMap, apiError := (&globalapi.MapServiceClient{Client: apiClient}).GetMapByName(mapName)
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
		return "Extreme"
	case 7:
		return "Death"
	}

	return "Unknown"
}
