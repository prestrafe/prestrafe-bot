package globalapi

import (
	"fmt"
	"strconv"

	"gitlab.com/prestrafe/prestrafe-bot/utils"
)

type Record struct {
	Id             int32   `json:"id"`
	SteamId64      string  `json:"steamid64"`
	PlayerName     string  `json:"player_name"`
	SteamId        string  `json:"steam_id"`
	ServerId       int32   `json:"server_id"`
	MapId          int32   `json:"map_id"`
	Stage          int32   `json:"stage"`
	Mode           string  `json:"mode"`
	TickRate       int32   `json:"tickrate"`
	Time           float32 `json:"time"`
	Teleports      int32   `json:"teleports"`
	CreatedOn      string  `json:"created_on"`
	UpdatedOn      string  `json:"updated_on"`
	UpdatedBy      int64   `json:"updated_by"`
	RecordFilterId int32   `json:"record_filter_id"`
	ServerName     string  `json:"server_name"`
	MapName        string  `json:"map_name"`
	Points         int32   `json:"points"`
	ReplayId       int32   `json:"replay_id"`
}

func (record *Record) FormattedTime() string {
	duration := int64(record.Time * 1_000)
	return fmt.Sprintf("%02d:%02d.%03d", duration/60_000, (duration%60_000)/1_000, duration%1_000)
}

type RecordServiceClient struct {
	Client
}

func (s *RecordServiceClient) GetRecordsTop(criteria QueryParameters) (result []Record, err error) {
	result = []Record{}
	err = s.GetWithParameters("records/top", criteria, &result)

	return
}

func (s *RecordServiceClient) GetWorldRecord(mapName, mode string, stage int) (nub, pro *Record, err error) {
	nubs, err := s.GetRecordsTop(QueryParameters{
		"map_name":          mapName,
		"modes_list_string": mode,
		"tickrate":          "128",
		"stage":             strconv.Itoa(stage),
		"overall":           "true",
		"limit":             "1",
	})
	if err != nil {
		return nil, nil, err
	}

	pros, err := s.GetRecordsTop(QueryParameters{
		"map_name":          mapName,
		"modes_list_string": mode,
		"has_teleports":     "false",
		"tickrate":          "128",
		"stage":             strconv.Itoa(stage),
		"overall":           "true",
		"limit":             "1",
	})
	if err != nil {
		return nil, nil, err
	}

	if len(nubs) > 0 {
		nub = &nubs[0]
	}
	if len(pros) > 0 {
		pro = &pros[0]
	}

	return
}

func (s *RecordServiceClient) GetPersonalRecord(mapName, mode string, stage int, steamId64 int64) (nub, pro *Record, err error) {
	nubs, err := s.GetRecordsTop(QueryParameters{
		"map_name":          mapName,
		"modes_list_string": mode,
		"steam_id":          utils.ConvertSteamId(steamId64),
		"tickrate":          "128",
		"stage":             strconv.Itoa(stage),
		"limit":             "1",
	})
	if err != nil {
		return nil, nil, err
	}

	pros, err := s.GetRecordsTop(QueryParameters{
		"map_name":          mapName,
		"modes_list_string": mode,
		"steam_id":          utils.ConvertSteamId(steamId64),
		"has_teleports":     "false",
		"tickrate":          "128",
		"stage":             strconv.Itoa(stage),
		"limit":             "1",
	})
	if err != nil {
		return nil, nil, err
	}

	if len(nubs) > 0 {
		nub = &nubs[0]
	}
	if len(pros) > 0 {
		pro = &pros[0]
	}

	return
}
