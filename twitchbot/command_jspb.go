package twitchbot

import (
	"fmt"

	"github.com/gempir/go-twitch-irc"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func CreateHandleJSCommand(jumpType, jumpName string) CommandHandler {
	return func(user twitch.User, strings []string) string {
		gameState, err := gsi.GetGameState()
		if err != nil {
			return "Could not retrieve current game state"
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

func HandleLJPBCommand(user twitch.User, strings []string) string {
	gameState, err := gsi.GetGameState()
	if err != nil {
		return "Could not retrieve current game state"
	}

	jumpStat, err := globalapi.GetJumpStatPersonalBest("longjump", gameState.Player.SteamId)

	if jumpStat != nil && err != nil {
		return fmt.Sprintf("%s Long Jump record of %s: %.04f units", gameState.Player.Clan, gameState.Player.Name, jumpStat.Distance)
	} else {
		return fmt.Sprintf("%s Long Jump record of %s: None", gameState.Player.Clan, gameState.Player.Name)
	}
}
