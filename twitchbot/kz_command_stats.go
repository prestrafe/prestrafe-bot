package twitchbot

import (
	"errors"
	"fmt"

	"gitlab.com/prestrafe/prestrafe-bot/utils"

	"gitlab.com/prestrafe/prestrafe-bot/gsi"
)

func NewStatsCommand(gsiClient gsi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("stats").
		WithHandler(createStatsHandler(gsiClient))
}

func createStatsHandler(gsiClient gsi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil || !gameState.IsKZGameState() {
			return "", errors.New("could not retrieve KZ game play")
		}

		return fmt.Sprintf("Stats page for %s: https://gokzstats.com/?name=%s&mode=%s", ctx.Channel(), utils.ConvertSteamId(gameState.Provider.SteamId), gameState.Player.TimerMode()), nil
	}
}
