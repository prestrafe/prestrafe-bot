package twitchbot

import (
	"errors"
	"fmt"

	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
	"gitlab.com/prestrafe/prestrafe-bot/helper"
	"gitlab.com/prestrafe/prestrafe-bot/smclient"
)

// This is still not fully accurate, this should require checking with the API to see if a filter exists for this map, mode and course.
func NewGlobalCheckCommand(gsiClient gsiclient.Client, smClient smclient.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("globalcheck").
		WithAlias("gc").
		WithHandler(createGCHandler(gsiClient, smClient))
}

func createGCHandler(gsiClient gsiclient.Client, smClient smclient.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil || !gsiclient.IsKZGameState(gameState) {
			return "", errors.New("could not retrieve KZ game play")
		}
		fullPlayerState, smError := smClient.GetPlayerInfo()
		if smError != nil {
			return "", errors.New("Could not retrieve KZ gameplay from game server.")
		}
		if !helper.CompareData(fullPlayerState, gameState) {
			return "", errors.New("Could not match KZ gameplay from game server.")
		}
		if fullPlayerState.ServerGlobal == -1 {
			return "", errors.New("Could not retrieve server global status.")
		}
		if (fullPlayerState.ServerGlobal == 0) || !fullPlayerState.KZData.Global {
			return fmt.Sprintf("No"), nil
		} else {
			return fmt.Sprintf("Yes"), nil
		}
	}
}
