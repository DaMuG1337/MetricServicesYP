package config

import "flag"

type ConfigServer struct {
	FlagRunAddr string
}

func ServerFlags() ConfigServer {
	var c ConfigServer
	flag.StringVar(&c.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	return c
}
