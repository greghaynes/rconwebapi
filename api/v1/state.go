package v1

// StatusResponse is a summary status for the server
type StatusResponse struct {
	Hostname string        `json:"hostname"`
	Version  string        `json:"version"`
	Map      string        `json:"map"`
	Players  PlayersObject `json:"players"`
}

type PlayersObject struct {
	HumanPlayers int `json:"human_players"`
	BotPlayers   int `json:"bot_players"`
}
