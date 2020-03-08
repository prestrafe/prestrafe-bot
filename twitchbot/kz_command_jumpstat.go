package twitchbot

import (
	"errors"
	"fmt"

	"gitlab.com/prestrafe/prestrafe-bot/globalapi"
	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
)

func NewJumpStatCommand(gsiClient gsiclient.Client, apiClient globalapi.Client, name, jumpType, jumpName string, maxDistance int, withoutBinds bool) ChatCommandBuilder {
	return NewChatCommandBuilder(name).
		WithAlias(fmt.Sprintf("%spb", name), jumpType).
		WithParameter("binds", false, "(bind|nobind)").
		WithHandler(func(ctx CommandContext) (message string, err error) {
			noBind := withoutBinds
			if bindsParam, present := ctx.Parameter("noBind"); present {
				noBind = bindsParam == "nobind"
			}

			gameState, gsiError := gsiClient.GetGameState()
			if gsiError != nil || !gsiclient.IsKZGameState(gameState) {
				return "", errors.New("could not retrieve KZ game play")
			}

			jumpStat, apiError := (&globalapi.JumpStatServiceClient{Client: apiClient}).GetJumpStatPersonalBest(jumpType, maxDistance, gameState.Provider.SteamId, noBind)
			if jumpStat != nil && apiError == nil {
				binds := "no noBind"
				if jumpStat.HasBinds() {
					binds = "with noBind"
				}

				return fmt.Sprintf("%s record of %s: %.04f units, %d strafes, %s", jumpName, ctx.Channel(), jumpStat.Distance, jumpStat.StrafeCount, binds), nil
			} else {
				return fmt.Sprintf("%s record of %s: None", jumpName, ctx.Channel()), nil
			}
		})
}
