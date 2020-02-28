package old

import (
	"fmt"

	"github.com/gempir/go-twitch-irc"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func HandleWRCommand(user twitch.User, parameters []string) string {
	gameState, err := gsi.GetGameState()
	if err != nil {
		return "Could not retrieve KZ gameplay"
	}

	nub, pro, err := globalapi.GetWorldRecord(gameState.Map.Name, gameState.Player.TimerMode(), 0)

	message := fmt.Sprintf("Global Records on %s [%s]: ", gameState.Map.Name, gameState.Player.Clan)
	if nub != nil && err == nil {
		message += fmt.Sprintf("NUB: %s (%d TP) by %s", nub.FormattedTime(), nub.Teleports, nub.PlayerName)
	} else {
		message += fmt.Sprintf("NUB: None")
	}

	message += ", "

	if pro != nil && err == nil {
		message += fmt.Sprintf("PRO: %s by %s", pro.FormattedTime(), pro.PlayerName)
	} else {
		message += fmt.Sprintf("PRO: None")
	}

	return message
}
