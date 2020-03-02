package twitchbot

import (
	"errors"
	"fmt"

	"prestrafe-bot/globalapi"
	"prestrafe-bot/gsi"
)

func NewRankCommand(gsiClient gsi.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("rank").
		WithAlias("points").
		WithParameter("type", false, "(all|nub|pro|tp)").
		WithHandler(createRankHandler(gsiClient))
}

func createRankHandler(gsiClient gsi.Client) ChatCommandHandler {
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
		if gsiError != nil || !gameState.IsKZGameState() {
			return "", errors.New("could not retrieve KZ game play")
		}

		rank, apiError := globalapi.GetPlayerRank(gameState.Player.TimerModeId(), gameState.Provider.SteamId, teleports)
		if rank == nil || apiError != nil {
			return fmt.Sprintf("Points for %s [%s]: Unknown", gameState.Player.Name, gameState.Player.TimerModeName()), nil
		} else {
			return fmt.Sprintf("Points for %s [%s]: %d points with %d finishes", gameState.Player.Name, gameState.Player.TimerModeName(), rank.Points, rank.Finishes), nil
		}
	}
}
