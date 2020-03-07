package globalapi

import (
	"strconv"

	"gitlab.com/prestrafe/prestrafe-bot/utils"
)

type JumpStat struct {
	Id            int     `json:"id"`
	ServerId      int     `json:"server_id"`
	SteamId64     int64   `json:"steamid64"`
	PlayerName    string  `json:"player_name"`
	SteamId       string  `json:"steam_id"`
	JumpType      int     `json:"jump_type"`
	Distance      float64 `json:"distance"`
	TickRate      int     `json:"tickrate"`
	MslCount      int     `json:"msl_count"`
	StrafeCount   int     `json:"strafe_count"`
	IsCrouchBind  int     `json:"is_crouch_bind"`
	IsForwardBind int     `json:"is_forward_bind"`
	IsCrouchBoost int     `json:"is_crouch_boost"`
	UpdatedById   int     `json:"updated_by_id"`
	CreatedOn     string  `json:"created_on"`
	UpdatedOn     string  `json:"updated_on"`
}

func (js *JumpStat) HasBinds() bool {
	return js.IsCrouchBind != 0 || js.IsForwardBind != 0
}

type JumpStatServiceClient struct {
	Client
}

func (s *JumpStatServiceClient) GetJumpStats(criteria QueryParameters) (result []JumpStat, err error) {
	result = []JumpStat{}
	err = s.GetWithParameters("jumpstats", criteria, &result)

	return
}

func (s *JumpStatServiceClient) GetJumpStatPersonalBest(jumpType string, maxDistance int, steamId64 int64) (jumpStat *JumpStat, err error) {
	jumpStats, err := s.GetJumpStats(QueryParameters{
		"jumptype":           jumpType,
		"steam_id":           utils.ConvertSteamId(steamId64),
		"less_than_distance": strconv.Itoa(maxDistance),
		"limit":              "25",
	})

	if jumpStats != nil && len(jumpStats) > 0 {
		for i, js := range jumpStats {
			if js.StrafeCount <= 15 {
				jumpStat = &jumpStats[i]
				return
			}
		}
	}

	return
}
