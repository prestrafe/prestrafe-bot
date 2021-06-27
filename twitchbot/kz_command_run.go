package twitchbot

import (
	"errors"
	"fmt"
	"math"

	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
	"gitlab.com/prestrafe/prestrafe-bot/helper"
	"gitlab.com/prestrafe/prestrafe-bot/smclient"
)

func NewRunCommand(gsiClient gsiclient.Client, smClient smclient.Client) ChatCommandBuilder {
	return NewChatCommandBuilder("run").
		WithHandler(createRunHandler(gsiClient, smClient))
}

func createRunHandler(gsiClient gsiclient.Client, smClient smclient.Client) ChatCommandHandler {
	return func(ctx CommandContext) (message string, err error) {
		gameState, gsiError := gsiClient.GetGameState()
		if gsiError != nil || !gsiclient.IsKZGameState(gameState) {
			return "", errors.New("could not retrieve KZ gameplay")
		}
		fullPlayerState, smError := smClient.GetPlayerInfo()
		if smError != nil {
			return "", errors.New("could not retrieve data from game server")
		}
		if !helper.CompareData(fullPlayerState, gameState) {
			return "", errors.New("could not match data from game server with GSI client")
		}

		var course string
		// Course
		if fullPlayerState.KZData.Course == 0 {
			course = "Main"
		} else {
			course = fmt.Sprintf("Bonus %d", fullPlayerState.KZData.Course)
		}

		// Time
		hours := math.Floor(fullPlayerState.KZData.Time / 60 / 60)
		minutes := math.Floor(fullPlayerState.KZData.Time/60) - hours*60
		seconds := fullPlayerState.KZData.Time - hours*3600 - minutes*60

		var time string
		if hours > 0 {
			time = fmt.Sprintf("%02d:%02d:%05.2f", int(hours), int(minutes), seconds)
		} else if minutes > 0 {
			time = fmt.Sprintf("%02d:%05.2f", int(minutes), seconds)
		} else {
			time = fmt.Sprintf("%05.2f", seconds)
		}

		return fmt.Sprintf("Map: %s, Course: %s, Checkpoints: %d, Teleports: %d, Time elapsed: %s", fullPlayerState.MapName, course, fullPlayerState.KZData.Checkpoints, fullPlayerState.KZData.Teleports, time), nil
	}
}
