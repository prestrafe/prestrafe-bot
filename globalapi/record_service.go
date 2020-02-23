package globalapi

import (
	"fmt"
)

type Record struct {
	Id             int32   `json:"id"`
	SteamId64      int64   `json:"steamid64"`
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

func GetRecordsTop(criteria QueryParameters) (result []Record, err error) {
	result = []Record{}
	err = globalApiGet("records/top", &result, criteria)

	return
}

func GetWorldRecord(mapName, mode string) (nub, pro *Record, err error) {
	nubs, err := GetRecordsTop(QueryParameters{
		"map_name":          mapName,
		"modes_list_string": mode,
		"has_teleports":     "true",
		"tickrate":          "128",
		"stage":             "0",
		"overall":           "true",
		"limit":             "1",
	})
	pros, err := GetRecordsTop(QueryParameters{
		"map_name":          mapName,
		"modes_list_string": mode,
		"has_teleports":     "false",
		"tickrate":          "128",
		"stage":             "0",
		"overall":           "true",
		"limit":             "1",
	})

	if nubs != nil && len(nubs) > 0 {
		nub = &nubs[0]
	}
	if pros != nil && len(pros) > 0 {
		pro = &pros[0]
	}

	return
}

func GetPersonalRecord(mapName string, mode string, steamId64 int64) (nub, pro *Record, err error) {
	nubs, err := GetRecordsTop(QueryParameters{
		"map_name":          mapName,
		"modes_list_string": mode,
		"steamid64":         string(steamId64),
		"has_teleports":     "true",
		"tickrate":          "128",
		"stage":             "0",
		"overall":           "true",
		"limit":             "1",
	})
	pros, err := GetRecordsTop(QueryParameters{
		"map_name":          mapName,
		"modes_list_string": mode,
		"steamid64":         string(steamId64),
		"has_teleports":     "false",
		"tickrate":          "128",
		"stage":             "0",
		"overall":           "true",
		"limit":             "1",
	})

	if nubs != nil && len(nubs) > 0 {
		nub = &nubs[0]
	}
	if pros != nil && len(pros) > 0 {
		pro = &pros[0]
	}

	return
}
