package gsiclient

import (
	"regexp"
	"strings"

	"gitlab.com/prestrafe/prestrafe-gsi"
)

func IsKZGameState(gameState *gsi.GameState) bool {
	if gameState.Player == nil || gameState.Map == nil {
		return false
	}

	matchString, err := regexp.MatchString("^(workshop/[0-9]+/)?(kz|kzpro|skz|vnl|xc)_.+$", gameState.Map.Name)
	return matchString && err == nil
}

func GetMapName(mapState *gsi.MapState) string {
	if strings.HasPrefix(mapState.Name, "workshop") {
		return mapState.Name[strings.LastIndex(mapState.Name, "/")+1:]
	}

	return mapState.Name
}
func TimerMode(player *gsi.PlayerState) string {
	switch player.Clan {
	case "SKZ":
		return "kz_simple"
	case "VNL":
		return "kz_vanilla"
	}

	return "kz_timer"
}

func TimerModeName(player *gsi.PlayerState) string {
	switch player.Clan {
	case "SKZ":
		return "SKZ"
	case "VNL":
		return "VNL"
	}

	return "KZT"
}

func TimerModeId(player *gsi.PlayerState) int {
	switch player.Clan {
	case "SKZ":
		return 201
	case "VNL":
		return 202
	}

	return 200
}
