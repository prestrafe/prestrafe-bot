package twitchbot

import (
	"errors"
	"fmt"

	"gitlab.com/prestrafe/prestrafe-bot/globalapi"
	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
)

func NewPBCommand(gsiClient gsiclient.Client, apiClient globalapi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("pb").
		WithAlias("pr").
		WithParameter("map", false, "[A-Za-z0-9_]+").
		WithParameter("mode", false, "(kzt|skz|vnl)").
		WithHandler(createPBHandler(gsiClient, apiClient))
}

func createPBHandler(gsiClient gsiclient.Client, apiClient globalapi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		mapName, hasMapName := ctx.Parameter("map")
		modeName, hasModeName := ctx.Parameter("mode")

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

		nub, pro, apiError := (&globalapi.RecordServiceClient{Client: apiClient}).GetPersonalRecord(mapName, gsiclient.TimerModeFromName(modeName), 0, gameState.Provider.SteamId)

		message = fmt.Sprintf("PB of %s on %s [%s]: ", ctx.Channel(), mapName, gsiclient.TimerModeName(gameState.Player))
		if nub != nil && apiError == nil {
			message += fmt.Sprintf("NUB: %s (%d TP)", nub.FormattedTime(), nub.Teleports)
		} else {
			message += fmt.Sprintf("NUB: None")
		}

		message += ", "

		if pro != nil && apiError == nil {
			message += fmt.Sprintf("PRO: %s", pro.FormattedTime())
		} else {
			message += fmt.Sprintf("PRO: None")
		}

		return
	}
}
