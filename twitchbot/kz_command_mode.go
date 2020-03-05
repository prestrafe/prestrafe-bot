package twitchbot

import (
	"errors"

	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
)

func NewModeCommand(gsiClient gsiclient.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("mode").
		WithHandler(createModeHandler(gsiClient))
}

func createModeHandler(gsiClient gsiclient.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil || !gsiclient.IsKZGameState(gameState) {
			return "", errors.New("could not retrieve KZ game play")
		}

		return "Currently playing on: " + gsiclient.TimerModeName(gameState.Player), nil
	}
}
