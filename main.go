package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	host := flag.String("host", "127.0.0.1:8099", "host to listen on")
	rconAddress := flag.String("rconAddress", "", "rcon server address for state API")
	rconPassword := flag.String("rconPassword", "", "rcon server password for state API")
	flag.Parse()

	cfg := &Config{
		BindAddress:  *host,
		RconAddress:  *rconAddress,
		RconPassword: *rconPassword,
	}

	r := mux.NewRouter()
	http.Handle("/", r)

	rconSrv := NewRconServer(cfg)
	rconSrv.SetupHandlers(r)
	log.Println("Started rcon server")

	if cfg.RconAddress != "" {
		stateSrv := NewStateServer(cfg)
		stateSrv.SetupHandlers(r)
		log.Println("Started state server")
	}
	log.Fatal(http.ListenAndServe(cfg.BindAddress, nil))
}
