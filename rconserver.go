package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	v1 "github.com/greghaynes/rconwebapi/api/v1"
)

// RconServer manages server state
type RconServer struct {
	config *Config
}

func invalidMethod(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed."))
}

// NewRconServer creates a new server
func NewRconServer(config *Config) *RconServer {
	return &RconServer{
		config: config,
	}
}

// SetupHandlers adds handlers for rcon server
func (s *RconServer) SetupHandlers(r *mux.Router) {
	r.HandleFunc("/", s.indexHandler).Methods("GET")

	// Deprecated unversioned URL == v1
	r.HandleFunc("/rcon", s.rconHandler).Methods("POST")
	r.HandleFunc("/rcon_ws", s.rconWSHandler).Methods("POST")

	// v1 handlers
	r.HandleFunc("/v1/rcon", s.rconHandler).Methods("POST")
	r.HandleFunc("/v1/rcon_ws", s.rconWSHandler).Methods("POST")
}

func (s *RconServer) indexHandler(w http.ResponseWriter, req *http.Request) {
	LogRequest(req)
	w.Write([]byte("Hello!"))
}

func (s *RconServer) rconWSHandler(w http.ResponseWriter, req *http.Request) {
	LogRequest(req)

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("Failed to upgrade request to websocket.")
		return
	}
	defer conn.Close()

	var rconClient *RconClient

	for {
		var req v1.RconWSRequest
		if err := conn.ReadJSON(&req); err != nil {
			log.Printf("Failed to read websocket message: %v\n", err)
			return
		}

		if req.RequestType == v1.WSRequestTypeConnect {
			if rconClient != nil {
				log.Println("Got connect request when already connected")
				continue
			}

			var connectReq v1.RconWSConnectRequest
			if err = json.Unmarshal(req.Request, &connectReq); err != nil {
				log.Printf("Failed to parse connect request: %v\n", err)
				continue
			}

			rconClient, err = NewRconClient(connectReq.Address, connectReq.Password)
			if err != nil {
				log.Printf("Failed to connect to rcon: %v\n", err)
			}
			defer rconClient.Close()
		} else if req.RequestType == v1.WSRequestTypeCommand {
			if rconClient == nil {
				log.Println("Got command request while unconnected")
				continue
			}

			var commandReq v1.RconWSCommandRequest
			if err = json.Unmarshal(req.Request, &commandReq); err != nil {
				log.Printf("Failed to parse command request: %v\n", err)
				continue
			}

			resp, err := rconClient.Execute(commandReq.Command)
			if err != nil {
				log.Printf("Error executing rcon command: %v\n", err)
				continue
			}

			commandResp, err := json.Marshal(v1.RconWsCommandResponse{
				Output: resp,
			})
			if err != nil {
				log.Printf("Failed to marshall command response: %v\n", err)
				continue
			}
			response, err := json.Marshal(v1.RconWSResponse{
				ResponseType: v1.WSResponseTypeCommand,
				Response:     commandResp,
			})
			if err != nil {
				log.Printf("Failed to marshall websocket response: %v\n", err)
				continue
			}

			if err = conn.WriteMessage(websocket.TextMessage, response); err != nil {
				log.Printf("Failed to send websocket message: %v\n", err)
			}
		}
	}
}

func (s *RconServer) rconHandler(w http.ResponseWriter, req *http.Request) {
	LogRequest(req)

	ct := req.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write(([]byte("Ivalid Content-Type, only application/json allowed.")))
		return
	}

	decoder := json.NewDecoder(req.Body)
	var reqBody v1.RconReqBody
	if err := decoder.Decode(&reqBody); err != nil {
		log.Printf("Failed to parse request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request, unable to parse request body."))
		return
	}

	resp, err := s.makeRconRequest(&reqBody.RconRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error, command failed."))
		return
	}
	respBody := v1.RconResponseBody{
		RconResponse: v1.RconResponse{
			Output: resp,
		},
	}

	js, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Failed to serialize response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error, failed to marshall response."))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *RconServer) makeRconRequest(rconReq *v1.RconRequest) (string, error) {
	return MakeRconRequest(rconReq.Address, rconReq.Password, rconReq.Command)
}
