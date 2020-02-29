package twitchbot

import (
	"errors"
	"fmt"
	"strconv"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func NewBPBCommand(gsiClient gsi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("bpb").
		WithAlias("bpr").
		WithParameter("bonus", false, "[0-9]").
		WithParameter("map", false, "[A-Za-z0-9_]+").
		WithHandler(createBPBHandler(gsiClient))
}

func createBPBHandler(gsiClient gsi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		bonus, hasBonus := ctx.Parameter("bonus")
		mapName, hasMapName := ctx.Parameter("map")

		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil {
			return "", errors.New("could not retrieve KZ game play")
		}

		if !hasBonus {
			bonus = "1"
		}
		if !hasMapName {
			mapName = gameState.Map.Name
		}

		bonusNumber, _ := strconv.Atoi(bonus)
		if bonusNumber < 1 {
			return fmt.Sprintf("'%s' is not a valid bonus number.", bonus), nil
		}

		nub, pro, apiError := globalapi.GetPersonalRecord(mapName, gameState.Player.TimerMode(), bonusNumber, gameState.Player.SteamId)

		message = fmt.Sprintf("PB of %s on %s Bonus %d [%s]: ", gameState.Player.Name, mapName, bonusNumber, gameState.Player.Clan)
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
