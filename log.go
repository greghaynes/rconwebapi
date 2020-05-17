package main

import (
	"log"
	"net/http"
)

func LogRequest(req *http.Request) {
	log.Printf("Got %q request from %q for %q\n", req.Method, req.RemoteAddr, req.URL)
}
