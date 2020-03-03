package twitchbot

import (
	"errors"

	"prestrafe-bot/gsi"
)

func NewModeCommand(gsiClient gsi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("mode").
		WithHandler(createModeHandler(gsiClient))
}

func createModeHandler(gsiClient gsi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil || !gameState.IsKZGameState() {
			return "", errors.New("could not retrieve KZ game play")
		}

		return "Currently playing on: " + gameState.Player.TimerModeName(), nil
	}
}
