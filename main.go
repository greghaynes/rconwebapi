package main

import "flag"

func main() {
	host := flag.String("host", "127.0.0.1:8099", "host to listen on")
	flag.Parse()

	cfg := &Config{
		BindAddress: *host,
	}
	srv := NewServer(cfg)
	srv.Run()
}
