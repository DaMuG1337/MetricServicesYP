package config

import (
	"flag"
	"os"
	"strconv"
)

type ConfigServer struct {
	FlagRunAddr string
	Save        SaveParam
}

type SaveParam struct {
	Interval int
	FilePath string
	LoadFile bool
}

func ServerFlags() ConfigServer {
	var c ConfigServer
	c.FlagRunAddr = "localhost:8080"
	c.Save.Interval = 300
	c.Save.FilePath = "metrics.json"
	c.Save.LoadFile = false

	flag.StringVar(&c.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&c.Save.Interval, "i", 300, "interval save to file")
	flag.StringVar(&c.Save.FilePath, "f", "metrics.json", "file name save to metrics")
	flag.BoolVar(&c.Save.LoadFile, "r", false, "load old metrics")
	flag.Parse()

	if v := os.Getenv("ADDRESS"); v != "" {
		c.FlagRunAddr = v
	}
	if v := os.Getenv("STORE_INTERVAL"); v != "" {
		if b, err := strconv.Atoi(v); err == nil {
			c.Save.Interval = b
		}
	}
	if v := os.Getenv("FILE_STORAGE_PATH"); v != "" {
		c.Save.FilePath = v
	}
	if v := os.Getenv("RESTORE"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			c.Save.LoadFile = b
		}
	}
	return c
}
