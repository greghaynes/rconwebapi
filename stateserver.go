package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	v1 "github.com/greghaynes/rconwebapi/api/v1"
)

const (
	rconCommandStatus = "status"
)

// StateServer manages state server state
type StateServer struct {
	config *Config
}

// NewStateServer creates a StateServer
func NewStateServer(cfg *Config) *StateServer {
	return &StateServer{
		config: cfg,
	}
}

// SetupHandlers adds handlers for StateServer
func (s *StateServer) SetupHandlers(r *mux.Router) {
	r.HandleFunc("/v1/state/status", s.handleStatus).Methods("GET")
}

func (s *StateServer) handleStatus(w http.ResponseWriter, req *http.Request) {
	LogRequest(req)

	status, err := s.makeRconRequest(rconCommandStatus)
	if err != nil {
		log.Printf("RCON status command failed: %v\n", err)
		return
	}

	resp := statusToResponse(status)
	out, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal status response: %v\n", err)
		return
	}
	w.Write(out)
}

func statusToResponse(status string) *v1.StatusResponse {
	statusSplit := strings.SplitN(status, "\n", 8)

	hostname := removeFieldPrefix(statusSplit[0])
	version := removeFieldPrefix(statusSplit[1])
	map_name := removeFieldPrefix(statusSplit[5])
	playersline := removeFieldPrefix(statusSplit[6])
	rest := statusSplit[7]

	humans, bots := parsePlayersLine(playersline)
	allPlayers := parseAllPlayers(rest)
	players := v1.PlayersObject{HumanPlayers: humans, BotPlayers: bots, Players: allPlayers}

	return &v1.StatusResponse{
		Hostname: hostname,
		Version:  version,
		Map:      map_name,
		Players:  players,
	}
}

func parsePlayersLine(line string) (int, int) {
	lineSplit := strings.SplitN(line, " ", 4)
	humans, err := strconv.Atoi(lineSplit[0])
	if err != nil {
		log.Fatal(err)
	}
	bots, err := strconv.Atoi(lineSplit[2])
	if err != nil {
		log.Fatal(err)
	}
	return humans, bots

}

func parseAllPlayers(rest string) []v1.Player {
	players := make([]v1.Player, 0)
	//clean up input string
	trimmed := strings.TrimPrefix(strings.TrimSuffix(rest, "\n#end\n"), "\n")
	lines := strings.Split(trimmed, "\n")
	for idx, line := range lines {
		// removes header line
		if idx == 0 {
			continue
		}
		split := strings.Split(line, " ")
		// skip bots
		if split[2] == "BOT" {
			continue
		}
		ping, err := strconv.Atoi(split[6])
		if err != nil {
			log.Fatal(err)
		}
		loss, err := strconv.Atoi(split[7])
		if err != nil {
			log.Fatal(err)
		}

		player := v1.Player{
			Userid:        split[1],
			Name:          split[3],
			Uniqueid:      split[4],
			TimeConnected: split[5],
			Ping:          ping,
			Loss:          loss,
			State:         split[8],
			Rate:          split[9],
			Address:       split[10],
		}
		players = append(players, player)

	}
	return players
}

func removeFieldPrefix(line string) string {
	lineSplit := strings.SplitN(line, ": ", 2)
	return lineSplit[1]
}

func (s *StateServer) makeRconRequest(command string) (string, error) {
	return MakeRconRequest(s.config.RconAddress, s.config.RconPassword, command)
}
