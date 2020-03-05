package twitchbot

import (
	"errors"
	"fmt"
	"strconv"

	"gitlab.com/prestrafe/prestrafe-bot/globalapi"
	"gitlab.com/prestrafe/prestrafe-bot/gsi"
)

func NewBWRCommand(gsiClient gsi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("bwr").
		WithAlias("bgr", "bgwr", "btop").
		WithParameter("bonus", false, "[0-9]").
		WithParameter("map", false, "[A-Za-z0-9_]+").
		WithHandler(createBWRHandler(gsiClient))
}

func createBWRHandler(gsiClient gsi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		bonus, hasBonus := ctx.Parameter("bonus")
		mapName, hasMapName := ctx.Parameter("map")

		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil || !gameState.IsKZGameState() {
			return "", errors.New("could not retrieve KZ game play")
		}

		if !hasBonus {
			bonus = "1"
		}
		if !hasMapName {
			mapName = gameState.Map.GetMapName()
		}

		bonusNumber, _ := strconv.Atoi(bonus)
		if bonusNumber < 1 {
			return fmt.Sprintf("'%s' is not a valid bonus number.", bonus), nil
		}

		nub, pro, apiError := globalapi.GetWorldRecord(mapName, gameState.Player.TimerMode(), bonusNumber)

		message = fmt.Sprintf("Global Records on %s Bonus %d [%s]: ", mapName, bonusNumber, gameState.Player.TimerModeName())
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
