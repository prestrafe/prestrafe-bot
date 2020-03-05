package twitchbot

import (
	"errors"
	"fmt"

	"gitlab.com/prestrafe/prestrafe-bot/utils"

	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
)

func NewStatsCommand(gsiClient gsiclient.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("stats").
		WithHandler(createStatsHandler(gsiClient))
}

func createStatsHandler(gsiClient gsiclient.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil || !gsiclient.IsKZGameState(gameState) {
			return "", errors.New("could not retrieve KZ game play")
		}

		return fmt.Sprintf("Stats page for %s: https://gokzstats.com/?name=%s&mode=%s", ctx.Channel(), utils.ConvertSteamId(gameState.Provider.SteamId), gsiclient.TimerMode(gameState.Player)), nil
	}
}
