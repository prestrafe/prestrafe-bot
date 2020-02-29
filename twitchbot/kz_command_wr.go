package twitchbot

import (
	"errors"
	"fmt"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func CreateWrCommandBuilder() ChatCommandBuilder {
	return NewChatCommandBuilder("wr").
		WithAlias("gr", "gwr", "top").
		WithParameter("map", false, "[A-Za-z0-9_]").
		WithHandler(createWrHandler(nil))
}

func createWrHandler(gsiClient gsi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		mapName, hasMapName := ctx.Parameter("map")

		gameState, gsiError := gsiClient.GetGameState(ctx.GsiToken())
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
}
