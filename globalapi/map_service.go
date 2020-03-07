package globalapi

import "strconv"

type KzMap struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	FileSize            int    `json:"filesize"`
	Validated           bool   `json:"validated"`
	Difficulty          int    `json:"difficulty"`
	CreatedOn           string `json:"created_on"`
	UpdatedOn           string `json:"updated_on"`
	ApprovedBySteamId64 int64  `json:"approved_by_steamid64"`
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
