package v1

import "encoding/json"

const (
	// WSRequestTypeConnect is the request_type for connect messages
	WSRequestTypeConnect = "connect"
	// WSRequestTypeCommand is the request_type for command messages
	WSRequestTypeCommand = "command"
	// WSResponseTypeConnect is the response_type for connect messages
	WSResponseTypeConnect = "connect"
	// WSResponseTypeCommand is the response type for command messages
	WSResponseTypeCommand = "command"
)

// RconRequest is used for the standard HTTP POST API
type RconRequest struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	Command  string `json:"command"`
}

// RconResponse is the response for standard HTTP POST API
type RconResponse struct {
	Output string `json:"output"`
}

// RconWSConnectRequest connects to the rcon server via websocket
type RconWSConnectRequest struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

// RconWSCommandRequest sends a command to rcon via websocket
type RconWSCommandRequest struct {
	Command string `json:"command"`
}

// RconWSRequest is a container for websocket messages
type RconWSRequest struct {
	RequestType string          `json:"request_type"`
	Request     json.RawMessage `json:"request"`
}

// RconWsCommandResponse is the response for a websocket command
type RconWsCommandResponse struct {
	Output string `json:"output"`
}

// RconWSResponse is the response for ws messages
type RconWSResponse struct {
	ResponseType string          `json:"response_type"`
	Response     json.RawMessage `json:"response"`
}

// RconReqBody is the container for HTTP POST messages
type RconReqBody struct {
	RconRequest RconRequest `json:"rcon_request"`
}

// RconResponseBody is the container for HTTP POST responses
type RconResponseBody struct {
	RconResponse RconResponse `json:"rcon_response"`
}
