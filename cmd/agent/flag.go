package main

import (
	"flag"
	"os"
)

var flagRunAddr string
var flagReportInterval int
var flagPollInterval int

func parseFlags() {
	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	CommandLine.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	CommandLine.IntVar(&flagReportInterval, "r", 10, "report interval in seconds")
	CommandLine.IntVar(&flagPollInterval, "p", 2, "poll interval in seconds")
	CommandLine.Parse(os.Args[1:])
}
