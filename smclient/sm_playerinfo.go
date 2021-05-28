package smclient

type FullPlayerInfo struct {
	TimeStamp      int     `json:"timestamp"`
	AuthKey        string  `json:"authkey"`
	TimeoutsCTPrev int     `json:"timeoutsCTprev"`
	TimeoutsTPrev  int     `json:"timeoutsTprev"`
	TimeoutsCT     int     `json:"timeoutsCT"`
	TimeoutsT      int     `json:"timeoutsT"`
	ServerName     string  `json:"servername"`
	MapName        string  `json:"mapname"`
	ServerGlobal   int     `json:"serverglobal"`
	SteamId        int64   `json:"steamid,string"`
	Clan           string  `json:"clan"`
	Name           string  `json:"name"`
	TimeInServer   float64 `json:"timeinserver"` // Need a better name
	KZData         KZData  `json:"KZData"`
}

type KZData struct {
	Global      bool    `json:"global"`
	Course      int     `json:"course"`
	Time        float64 `json:"time"`
	Checkpoints int     `json:"checkpoints"`
	Teleports   int     `json:"teleports"`
}
