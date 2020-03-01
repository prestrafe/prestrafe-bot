package gsi

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

func (gameState *GameState) IsKZGameState() bool {
	if gameState.Player == nil || gameState.Map == nil {
		return false
	}

	matchString, err := regexp.MatchString("^(workshop/[0-9]+/)?(kz|kzpro|skz|vnl|xc)_.+$", gameState.Map.Name)
	return matchString && err == nil
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

func (mapState *MapState) GetMapName() string {
	if strings.HasPrefix(mapState.Name, "workshop") {
		return mapState.Name[strings.LastIndex(mapState.Name, "/")+1:]
	}

	return mapState.Name
}

type PlayerState struct {
	SteamId    int64       `json:"steamid,string"`
	Clan       string      `json:"clan"`
	Name       string      `json:"name"`
	MatchStats *MatchStats `json:"match_stats"`
}

func (player *PlayerState) TimerMode() string {
	switch player.Clan {
	case "SKZ":
		return "kz_simple"
	case "VNL":
		return "kz_vanilla"
	}

	return "kz_timer"
}

type MatchStats struct {
	Kills   int `json:"kills"`
	Assists int `json:"assists"`
	Deaths  int `json:"deaths"`
	Mvps    int `json:"mvps"`
	Score   int `json:"score"`
}
