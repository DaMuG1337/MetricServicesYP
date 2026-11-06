package config

import (
	"flag"
	"os"
	"strconv"
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
	c.Server.FlagRunAddr = "localhost:8080"
	c.Poll.FlagReportInterval = 10
	c.Poll.FlagPollInterval = 2

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	CommandLine.StringVar(&c.Server.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	CommandLine.IntVar(&c.Poll.FlagReportInterval, "r", 10, "report interval in seconds")
	CommandLine.IntVar(&c.Poll.FlagPollInterval, "p", 2, "poll interval in seconds")
	CommandLine.Parse(os.Args[1:])

	if v := os.Getenv("ADDRESS"); v != "" {
		c.Server.FlagRunAddr = v
	}
	if v := os.Getenv("REPORT_INTERVAL"); v != "" {
		c.Poll.FlagReportInterval, _ = strconv.Atoi(v)
	}
	if v := os.Getenv("POLL_INTERVAL"); v != "" {
		c.Poll.FlagPollInterval, _ = strconv.Atoi(v)
	}
	return c
}
