package globalapi

import "strconv"

type PlayerRank struct {
	Points     int     `json:"points"`
	Average    float64 `json:"average"`
	Rating     float64 `json:"rating"`
	Finishes   int     `json:"finishes"`
	SteamId64  string  `json:"steamid64,string"`
	SteamId    string  `json:"steamid"`
	PlayerName string  `json:"player_name"`
}

type PlayerRankServiceClient struct {
	Client
}

func (s *PlayerRankServiceClient) GetPlayerRanks(criteria QueryParameters) (result []PlayerRank, err error) {
	result = []PlayerRank{}
	err = s.GetWithParameters("player_ranks", criteria, &result)

	return
}

func (s *PlayerRankServiceClient) GetPlayerRank(modeId int, steamId64 int64, hasTeleports *bool) (rank *PlayerRank, err error) {
	teleports := ""
	if hasTeleports != nil {
		if *hasTeleports == true {
			teleports = "true"
		} else {
			teleports = "false"
		}
	}

	ranks, err := s.GetPlayerRanks(QueryParameters{
		"steamid64s":            strconv.FormatInt(steamId64, 10),
		"has_teleports":         teleports,
		"mode_ids":              strconv.Itoa(modeId),
		"stages":                "0",
		"finishes_greater_than": "0",
		"limit":                 "1",
	})
	if ranks != nil && len(ranks) > 0 {
		rank = &ranks[0]
	}

	return
}
