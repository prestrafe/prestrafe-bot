package twitchbot

import (
	"errors"
	"fmt"

	"gitlab.com/prestrafe/prestrafe-bot/globalapi"
	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
)

func NewRankCommand(gsiClient gsiclient.Client, apiClient globalapi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("rank").
		WithAlias("points").
		WithParameter("type", false, "(all|nub|pro|tp)").
		WithHandler(createRankHandler(gsiClient, apiClient))
}

func createRankHandler(gsiClient gsiclient.Client, apiClient globalapi.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		_type, hasType := ctx.Parameter("type")
		hasTeleports := false
		teleports := &hasTeleports

		if hasType && _type == "pro" {
			hasTeleports = false
		} else if hasType && _type == "tp" {
			hasTeleports = true
		} else {
			teleports = nil
		}

		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil || !gsiclient.IsKZGameState(gameState) {
			return "", errors.New("could not retrieve KZ game play")
		}

		rank, apiError := (&globalapi.PlayerRankServiceClient{Client: apiClient}).GetPlayerRank(gsiclient.TimerModeId(gameState.Player), gameState.Provider.SteamId, teleports)
		if rank == nil || apiError != nil {
			return fmt.Sprintf("Points for %s [%s]: Unknown", ctx.Channel(), gsiclient.TimerModeName(gameState.Player)), nil
		} else {
			return fmt.Sprintf("Points for %s [%s]: %d points with %d finishes", ctx.Channel(), gsiclient.TimerModeName(gameState.Player), rank.Points, rank.Finishes), nil
		}
	}
}
