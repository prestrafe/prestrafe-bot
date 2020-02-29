package twitchbot

import (
	"errors"
	"fmt"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func NewPBCommand(gsiClient gsi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("pb").
		WithAlias("pr").
		WithParameter("map", false, "[A-Za-z0-9_]+").
		WithHandler(createPBHandler(gsiClient))
}

func createPBHandler(gsiClient gsi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		mapName, hasMapName := ctx.Parameter("map")

		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil || !gameState.IsKZGameState() {
			return "", errors.New("could not retrieve KZ game play")
		}

		if !hasMapName {
			mapName = gameState.Map.GetMapName()
		}

		nub, pro, apiError := globalapi.GetPersonalRecord(mapName, gameState.Player.TimerMode(), 0, gameState.Player.SteamId)

		message = fmt.Sprintf("PB of %s on %s [%s]: ", gameState.Player.Name, gameState.Map.GetMapName(), gameState.Player.Clan)
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
