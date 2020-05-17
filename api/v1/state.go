package v1

// StatusResponse is a summary status for the server
type StatusResponse struct {
	Hostname string `json:"hostname"`
	Version  string `json:"version"`
}
