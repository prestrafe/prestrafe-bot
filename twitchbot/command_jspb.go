package twitchbot

import (
	"fmt"

	"github.com/gempir/go-twitch-irc"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func CreateJSHandler(jumpType, jumpName string) CommandHandler {
	return func(user twitch.User, parameters []string) string {
		gameState, err := gsi.GetGameState()
		if err != nil {
			return "Could not retrieve KZ gameplay"
		}

		jumpStat, err := globalapi.GetJumpStatPersonalBest(jumpType, gameState.Player.SteamId)

		if jumpStat != nil && err != nil {
			binds := "no binds"
			if jumpStat.HasBinds() {
				binds = "binds"
			}

			return fmt.Sprintf("%s record of %s: %.04f units, %d strafes, %s", jumpName, gameState.Player.Name, jumpStat.Distance, jumpStat.StrafeCount, binds)
		} else {
			return fmt.Sprintf("%s record of %s: None", jumpName, gameState.Player.Name)
		}
	}
}
