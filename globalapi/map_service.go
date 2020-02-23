package globalapi

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

func GetMaps(criteria QueryParameters) (result []KzMap, err error) {
	result = []KzMap{}
	err = globalApiGet("maps", &result, criteria)

	return
}

func GetMapById(id int) (result *KzMap, err error) {
	result = &KzMap{}
	err = globalApiGet("maps/"+string(id), result, nil)

	return
}

func GetMapByName(mapName string) (result *KzMap, err error) {
	result = &KzMap{}
	err = globalApiGet("maps/name/"+mapName, result, nil)

	return
}
