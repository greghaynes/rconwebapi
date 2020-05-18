package v1

// StatusResponse is a summary status for the server
type StatusResponse struct {
	Hostname string        `json:"hostname"`
	Version  string        `json:"version"`
	Map      string        `json:"map"`
	Players  PlayersObject `json:"players"`
}

type PlayersObject struct {
	HumanPlayers int      `json:"human_players"`
	BotPlayers   int      `json:"bot_players"`
	Players      []Player `json:"players"`
}

type Player struct {
	Userid        string `json:"user_id"`
	Name          string `json:"name"`
	Uniqueid      string `json:"unique_id"`
	TimeConnected string `json:"time_connected"`
	Ping          int    `json:"ping"`
	Loss          int    `json:"loss"`
	State         string `json:"state"`
	Rate          string `json:"rate"`
	Address       string `json:"address"`
}
