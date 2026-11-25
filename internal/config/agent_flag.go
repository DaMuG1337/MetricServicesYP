package config

import (
	"flag"
	"os"
	"strconv"

	"go.uber.org/zap"
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

func AgentFlags(logger *zap.SugaredLogger) ConfigAgent {
	var c ConfigAgent
	c.Server.FlagRunAddr = "localhost:8080"
	c.Poll.FlagReportInterval = 10
	c.Poll.FlagPollInterval = 2

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	CommandLine.StringVar(&c.Server.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	CommandLine.IntVar(&c.Poll.FlagReportInterval, "r", 10, "report interval in seconds")
	CommandLine.IntVar(&c.Poll.FlagPollInterval, "p", 2, "poll interval in seconds")
	CommandLine.Parse(os.Args[1:])

	if v, ok := os.LookupEnv("ADDRESS"); ok {
		c.Server.FlagRunAddr = v
	}
	if v, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		if n, err := strconv.Atoi(v); err == nil {
			c.Poll.FlagReportInterval = n
		} else {
			logger.Infof("невалидный REPORT_INTERVAL=%q: %v", v, err)
		}
	}
	if v, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		if n, err := strconv.Atoi(v); err == nil {
			c.Poll.FlagPollInterval = n
		} else {
			logger.Infof("невалидный POLL_INTERVAL=%q: %v", v, err)
		}
	}
	return c
}
