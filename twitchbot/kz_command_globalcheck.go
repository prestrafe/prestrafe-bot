package twitchbot

import (
	"errors"

	"gitlab.com/prestrafe/prestrafe-bot/globalapi"
	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
	"gitlab.com/prestrafe/prestrafe-bot/helper"
	"gitlab.com/prestrafe/prestrafe-bot/smclient"
)

func NewGlobalCheckCommand(gsiClient gsiclient.Client, smClient smclient.Client, apiClient globalapi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("globalcheck").
		WithAlias("gc").
		WithHandler(createGCHandler(gsiClient, smClient, apiClient))
}

func createGCHandler(gsiClient gsiclient.Client, smClient smclient.Client, apiClient globalapi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil || !gsiclient.IsKZGameState(gameState) {
			return "", errors.New("could not retrieve KZ game play")
		}
		fullPlayerState, smError := smClient.GetPlayerInfo()

		if smError != nil {
			return "", errors.New("could not retrieve KZ gameplay from game server")
		}
		if !helper.CompareData(fullPlayerState, gameState) {
			return "", errors.New("could not match KZ gameplay from game server")
		}
		if fullPlayerState.ServerGlobal == -1 {
			return "", errors.New("could not retrieve server global status")
		}
		if fullPlayerState.ServerGlobal == 0 {
			return "no (Server is not global)", nil
		}
		if !fullPlayerState.KZData.Global {
			return "no (Player is not verified)", nil
		}

		return (&globalapi.MapServiceClient{Client: apiClient}).CheckRecordFilter(fullPlayerState.KZData.Course, fullPlayerState.MapName, gsiclient.TimerModeId(gameState.Player)), nil
	}
}
