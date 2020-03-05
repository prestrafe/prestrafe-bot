package twitchbot

import (
	"errors"
	"fmt"

	"gitlab.com/prestrafe/prestrafe-bot/globalapi"
	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
)

func NewWRCommand(gsiClient gsiclient.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("wr").
		WithAlias("gr", "gwr", "top").
		WithParameter("map", false, "[A-Za-z0-9_]+").
		WithHandler(createWRHandler(gsiClient))
}

func createWRHandler(gsiClient gsiclient.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		mapName, hasMapName := ctx.Parameter("map")

		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil || !gsiclient.IsKZGameState(gameState) {
			return "", errors.New("could not retrieve KZ game play")
		}

		if !hasMapName {
			mapName = gsiclient.GetMapName(gameState.Map)
		}

		nub, pro, apiError := globalapi.GetWorldRecord(mapName, gsiclient.TimerMode(gameState.Player), 0)

		message = fmt.Sprintf("Global Records on %s [%s]: ", mapName, gsiclient.TimerModeName(gameState.Player))
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
}
