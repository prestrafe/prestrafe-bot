package gsiclient

import (
	"regexp"
	"strings"
)

type GameState struct {
	Auth     *AuthState     `json:"auth"`
	Map      *MapState      `json:"map"`
	Player   *PlayerState   `json:"player"`
	Provider *ProviderState `json:"provider"`
}

type AuthState struct {
	Token string `json:"token"`
}

type ProviderState struct {
	Name      string `json:"name"`
	AppId     int    `json:"appid"`
	Version   int    `json:"version"`
	SteamId   int64  `json:"steamid,string"`
	Timestamp int64  `json:"timestamp"`
}

type MapState struct {
	Name string `json:"name"`
}

type PlayerState struct {
	SteamId    int64       `json:"steamid,string"`
	Clan       string      `json:"clan"`
	Name       string      `json:"name"`
	MatchStats *MatchStats `json:"match_stats"`
}

type MatchStats struct {
	Kills   int `json:"kills"`
	Assists int `json:"assists"`
	Deaths  int `json:"deaths"`
	Mvps    int `json:"mvps"`
	Score   int `json:"score"`
}

func IsKZGameState(gameState *GameState) bool {
	if gameState.Player == nil || gameState.Map == nil {
		return false
	}

	mapName := strings.ToLower(gameState.Map.Name)
	matchString, err := regexp.MatchString("^(workshop/[0-9]+/)?(bkz|kz|kzpro|skz|vnl|xc)_.+$", mapName)
	return matchString && err == nil
}

func GetMapName(mapState *MapState) string {
	mapName := strings.ToLower(mapState.Name)
	if strings.HasPrefix(mapName, "workshop") {
		return mapName[strings.LastIndex(mapName, "/")+1:]
	}

	return mapName
}
func TimerMode(player *PlayerState) string {
	return TimerModeFromName(player.Clan)
}

func TimerModeFromName(name string) string {
	switch strings.ToLower(name) {
	case "skz":
		return "kz_simple"
	case "vnl":
		return "kz_vanilla"
	}

	return "kz_timer"
}

func TimerModeName(player *PlayerState) string {
	return TimerModeNameFromName(player.Clan)
}

func TimerModeNameFromName(name string) string {
	switch strings.ToUpper(name) {
	case "SKZ":
		return "SKZ"
	case "VNL":
		return "VNL"
	}

	return "KZT"
}

func TimerModeId(player *PlayerState) int {
	switch player.Clan {
	case "SKZ":
		return 201
	case "VNL":
		return 202
	}

	return 200
}
