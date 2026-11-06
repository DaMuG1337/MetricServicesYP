package config

import (
	"flag"
	"os"
)

type ConfigServer struct {
	FlagRunAddr string
}

func ServerFlags() ConfigServer {
	var c ConfigServer
	c.FlagRunAddr = "localhost:8080"

	flag.StringVar(&c.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if v := os.Getenv("ADDRESS"); v != "" {
		c.FlagRunAddr = v
	}
	return c
}
