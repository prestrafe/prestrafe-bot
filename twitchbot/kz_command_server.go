package twitchbot

import (
	"errors"
	"fmt"

	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
	"gitlab.com/prestrafe/prestrafe-bot/helper"
	"gitlab.com/prestrafe/prestrafe-bot/smclient"
)

func NewServerCommand(gsiClient gsiclient.Client, smClient smclient.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("server").
		WithHandler(createServerHandler(gsiClient, smClient))
}

func createServerHandler(gsiClient gsiclient.Client, smClient smclient.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil || !gsiclient.IsKZGameState(gameState) {
			return "", errors.New("Could not retrieve KZ gameplay.")
		}
		fullPlayerState, smError := smClient.GetPlayerInfo()
		if smError != nil {
			return "", errors.New("Could not retrieve data from game server.")
		}
		if !helper.CompareData(fullPlayerState, gameState) {
			return "", errors.New("Could not match data from game server with GSI client.")
		}

		global := "N/A"
		if fullPlayerState.ServerGlobal == 1 {
			global = "Yes"
		} else if fullPlayerState.ServerGlobal == 0 {
			global = "No"
		}
		return fmt.Sprintf("Current server: %s. Global status: %s", fullPlayerState.ServerName, global), nil
	}
}
