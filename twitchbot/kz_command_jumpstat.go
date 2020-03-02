package twitchbot

import (
	"errors"
	"fmt"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func NewJumpStatCommand(gsiClient gsi.Client, name, jumpType, jumpName string, maxDistance int) ChatCommandBuilder {
	return NewChatCommandBuilder(name).
		WithAlias(fmt.Sprintf("%spb", name), jumpType).
		WithHandler(func(ctx CommandContext) (message string, err error) {
			gameState, gsiError := gsiClient.GetGameState()
			if gsiError != nil || !gameState.IsKZGameState() {
				return "", errors.New("could not retrieve KZ game play")
			}

			jumpStat, apiError := globalapi.GetJumpStatPersonalBest(jumpType, maxDistance, gameState.Provider.SteamId)
			if jumpStat != nil && apiError == nil {
				binds := "no binds"
				if jumpStat.HasBinds() {
					binds = "with binds"
				}

				return fmt.Sprintf("%s record of %s: %.04f units, %d strafes, %s", jumpName, ctx.Channel(), jumpStat.Distance, jumpStat.StrafeCount, binds), nil
			} else {
				return fmt.Sprintf("%s record of %s: None", jumpName, ctx.Channel()), nil
			}
		})
}
