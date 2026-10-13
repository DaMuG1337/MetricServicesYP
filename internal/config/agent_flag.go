package config

import (
	"flag"
	"os"
)

type ConfigAgent struct {
	Server ServerConfig
	Poll   PollConfig
}

type ServerConfig struct {
	FlagRunAddr string
}

type PollConfig struct {
	FlagReportInterval int
	FlagPollInterval   int
}

func AgentFlags() ConfigAgent {
	var c ConfigAgent
	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	CommandLine.StringVar(&c.Server.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	CommandLine.IntVar(&c.Poll.FlagReportInterval, "r", 10, "report interval in seconds")
	CommandLine.IntVar(&c.Poll.FlagPollInterval, "p", 2, "poll interval in seconds")
	CommandLine.Parse(os.Args[1:])
	return c
}
