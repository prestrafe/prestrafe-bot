package globalapi

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

func GetJumpStatTop(jumpType string, criteria QueryParameters) (result []JumpStat, err error) {
	result = []JumpStat{}
	err = globalApiGet("jumpstats/"+jumpType+"/top", &result, criteria)

	return
}

func GetJumpStatPersonalBest(jumpType string, steamId64 int64) (jumpStat *JumpStat, err error) {
	jumpStats, err := GetJumpStatTop(jumpType, QueryParameters{"steam_id": convertSteamId(steamId64), "limit": "1"})
	if jumpStats != nil && len(jumpStats) > 0 {
		jumpStat = &jumpStats[0]
	}

	return
}
