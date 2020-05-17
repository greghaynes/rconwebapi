package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	WSRequestTypeConnect  = "connect"
	WSRequestTypeCommand  = "command"
	WSResponseTypeConnect = "connect"
	WSResponseTypeCommand = "command"
)

// Config holds configuration for the server
type Config struct {
	BindAddress string
}

type rconRequest struct {
	Address  string
	Password string
	Command  string
}

type rconResponse struct {
	Output string
}

type rconWSConnectRequest struct {
	Address  string
	Password string
}

type rconWSCommandRequest struct {
	Command string
}

type rconWSRequest struct {
	RequestType string
	Request     json.RawMessage
}

type rconWsCommandResponse struct {
	Output string
}

type rconWSResponse struct {
	ResponseType string
	Response     json.RawMessage
}

type rconReqBody struct {
	RconRequest rconRequest
}

type rconResponseBody struct {
	RconResponse rconResponse
}

// Server manages server state
type Server struct {
	config *Config
}

func logRequest(req *http.Request) {
	log.Printf("Got %q request from %q for %q\n", req.Method, req.RemoteAddr, req.URL)
}

func invalidMethod(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed."))
}

// NewServer creates a new server
func NewServer(config *Config) *Server {
	return &Server{
		config: config,
	}
}

// Run starts the Server
func (s *Server) Run() {
	s.setupHandlers()
	log.Fatal(http.ListenAndServe(s.config.BindAddress, nil))
}

func (s *Server) setupHandlers() {
	http.HandleFunc("/", s.indexHandler)
	http.HandleFunc("/rcon", s.rconHandler)
	http.HandleFunc("/rcon_ws", s.rconWSHandler)
}

func (s *Server) indexHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)

	if req.Method != http.MethodGet {
		invalidMethod(w)
		return
	}

	w.Write([]byte("Hello!"))
}

func (s *Server) rconWSHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)

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
		var req rconWSRequest
		if err := conn.ReadJSON(&req); err != nil {
			log.Printf("Failed to read websocket message: %v\n", err)
			return
		}

		if req.RequestType == WSRequestTypeConnect {
			if rconClient != nil {
				log.Println("Got connect request when already connected")
				continue
			}

			var connectReq rconWSConnectRequest
			if err = json.Unmarshal(req.Request, &connectReq); err != nil {
				log.Printf("Failed to parse connect request: %v\n", err)
				continue
			}

			rconClient, err = NewRconClient(connectReq.Address, connectReq.Password)
			if err != nil {
				log.Printf("Failed to connect to rcon: %v\n", err)
			}
			defer rconClient.Close()
		} else if req.RequestType == WSRequestTypeCommand {
			if rconClient == nil {
				log.Println("Got command request while unconnected")
				continue
			}

			var commandReq rconWSCommandRequest
			if err = json.Unmarshal(req.Request, &commandReq); err != nil {
				log.Printf("Failed to parse command request: %v\n", err)
				continue
			}

			resp, err := rconClient.Execute(commandReq.Command)
			if err != nil {
				log.Printf("Error executing rcon command: %v\n", err)
				continue
			}

			commandResp, err := json.Marshal(rconWsCommandResponse{
				Output: resp,
			})
			if err != nil {
				log.Printf("Failed to marshall command response: %v\n", err)
				continue
			}
			response, err := json.Marshal(rconWSResponse{
				ResponseType: WSResponseTypeCommand,
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

func (s *Server) rconHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)

	if req.Method != http.MethodPost {
		invalidMethod(w)
		return
	}

	ct := req.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write(([]byte("Ivalid Content-Type, only application/json allowed.")))
		return
	}

	decoder := json.NewDecoder(req.Body)
	var reqBody rconReqBody
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
	respBody := rconResponseBody{
		RconResponse: rconResponse{
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

func (s *Server) makeRconRequest(rconReq *rconRequest) (string, error) {
	client, err := NewRconClient(rconReq.Address, rconReq.Password)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Execute(rconReq.Command)
	if err != nil {
		return "", err
	}
	return resp, nil
}
