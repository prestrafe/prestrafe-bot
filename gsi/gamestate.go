package gsi

type GameState struct {
	Map    *MapState    `json:"map"`
	Player *PlayerState `json:"player"`
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

func (player *PlayerState) TimerMode() string {
	switch player.Clan {
	case "SKZ":
		return "kz_simple"
	case "VNL":
		return "kz_vanilla"
	}

	return "kz_timer"
}
