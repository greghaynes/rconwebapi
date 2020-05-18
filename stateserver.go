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

	humans, bots := parsePlayersLine(playersline)
	players := v1.PlayersObject{HumanPlayers: humans, BotPlayers: bots}

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

func removeFieldPrefix(line string) string {
	lineSplit := strings.SplitN(line, ": ", 2)
	return lineSplit[1]
}

func (s *StateServer) makeRconRequest(command string) (string, error) {
	return MakeRconRequest(s.config.RconAddress, s.config.RconPassword, command)
}
