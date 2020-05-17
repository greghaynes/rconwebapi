package main

import (
	"errors"

	"github.com/gorcon/rcon"
)

// RconClient acts as a client to an rcon server
type RconClient struct {
	address  string
	password string
	conn     *rcon.Conn
}

// NewRconClient returns an RconClient
func NewRconClient(address string, password string) (*RconClient, error) {
	client := &RconClient{
		address:  address,
		password: password,
	}

	if err := client.connect(); err != nil {
		return nil, err
	}
	return client, nil
}

// Execute runs a command on rcon server
func (r *RconClient) Execute(command string) (string, error) {
	if r.conn == nil {
		return "", errors.New("Not connected")
	}
	return r.conn.Execute(command)
}

// Close closes the connection
func (r *RconClient) Close() {
	r.conn.Close()
}

func (r *RconClient) connect() error {
	conn, err := rcon.Dial(r.address, r.password)
	if err != nil {
		return err
	}

	r.conn = conn
	return nil
}
