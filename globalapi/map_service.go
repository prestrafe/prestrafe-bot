package globalapi

import (
	"regexp"
	"strconv"
)

type KzMap struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	FileSize            int    `json:"filesize"`
	Validated           bool   `json:"validated"`
	Difficulty          int    `json:"difficulty"`
	CreatedOn           string `json:"created_on"`
	UpdatedOn           string `json:"updated_on"`
	ApprovedBySteamId64 string `json:"approved_by_steamid64"`
	WorkshopUrl         string `json:"workshop_url"`
	DownloadUrl         string `json:"download_url"`
}

type RecordFilter struct {
	Id        int32  `json:"id"`
	MapId     int32  `json:"map_id"`
	Stage     int32  `json:"stage"`
	Mode      string `json:"mode"`
	TickRate  int32  `json:"tickrate"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
	UpdatedBy int64  `json:"updated_by"`
}

type MapServiceClient struct {
	Client
}

func (s *MapServiceClient) GetMaps(criteria QueryParameters) (result []KzMap, err error) {
	result = []KzMap{}
	err = s.GetWithParameters("maps", criteria, &result)

	return
}

func (s *MapServiceClient) GetMapById(id int) (result *KzMap, err error) {
	result = &KzMap{}
	err = s.Get("maps/"+strconv.Itoa(id), result)

	return
}

func (s *MapServiceClient) GetMapByName(mapName string) (result *KzMap, err error) {
	result = &KzMap{}
	err = s.Get("maps/name/"+mapName, result)

	return
}

func (s *MapServiceClient) GetMapIdByName(mapName string) (result int, err error) {
	kzMap := &KzMap{}
	err = s.Get("maps/name/"+mapName, kzMap)

	return kzMap.Id, nil
}

func (s *MapServiceClient) CheckRecordFilter(stage int, mapName string, modeId int) string {
	mapId, err := s.GetMapIdByName(mapName)
	if err != nil {
		return "Cannot establish a connection to the API server."
	}
	if mapId != 0 {
		result := []RecordFilter{}

		err = s.GetWithParameters("record_filters", QueryParameters{
			"stages":   strconv.Itoa(stage),
			"map_ids":  strconv.Itoa(mapId),
			"mode_ids": strconv.Itoa(modeId),
		}, &result)
		if err != nil {
			match, _ := regexp.MatchString(`Expected \d+, but got \d+ instead!`, err.Error())
			if match {
				return "Cannot establish a connection to the API server."
			}
		} else if len(result) == 0 {
			return "No (Filter does not exist for this course)"
		} else {
			return "Yes"
		}
	}
	return "No (Map does not exist in API database)"
}
