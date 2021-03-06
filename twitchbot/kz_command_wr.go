package twitchbot

import (
	"errors"
	"fmt"

	"gitlab.com/prestrafe/prestrafe-bot/globalapi"
	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
)

func NewWRCommand(gsiClient gsiclient.Client, apiClient globalapi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("wr").
		WithAlias("gr", "gwr", "top").
		WithParameter("map", false, mapRegexPattern).
		WithParameter("mode", false, modeRegexPattern).
		WithHandler(createWRHandler(gsiClient, apiClient))
}

func createWRHandler(gsiClient gsiclient.Client, apiClient globalapi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		mapName, hasMapName := ctx.Parameter("map")
		modeName, hasModeName := ctx.Parameter("mode")

		if !hasMapName || !hasModeName {
			gameState, gsiError := gsiClient.GetGameState()
			if gsiError != nil || !gsiclient.IsKZGameState(gameState) {
				return "", errors.New("could not retrieve KZ game play")
			}

			if !hasMapName {
				mapName = gsiclient.GetMapName(gameState.Map)
			}
			if !hasModeName {
				modeName = gsiclient.TimerModeName(gameState.Player)
			}
		}

		nub, pro, apiError := (&globalapi.RecordServiceClient{Client: apiClient}).GetWorldRecord(mapName, gsiclient.TimerModeFromName(modeName), 0)

		message = fmt.Sprintf("Global Records on %s [%s]: ", mapName, gsiclient.TimerModeNameFromName(modeName))
		if nub != nil && apiError == nil {
			message += fmt.Sprintf("NUB: %s (%d TP) by %s", nub.FormattedTime(), nub.Teleports, nub.PlayerName)
		} else {
			message += "NUB: None"
		}

		message += ", "

		if pro != nil && apiError == nil {
			message += fmt.Sprintf("PRO: %s by %s", pro.FormattedTime(), pro.PlayerName)
		} else {
			message += "PRO: None"
		}

		return
	}
}
