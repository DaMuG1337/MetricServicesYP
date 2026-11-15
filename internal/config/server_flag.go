package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/DaMuG1337/MetricServices/internal/logger"
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
	logger, err := logger.New()
	if err != nil {
		logger.Error("Ошибка запуска логгера", err)
	}
	c.FlagRunAddr = "localhost:8080"
	c.Save.Interval = 300
	c.Save.FilePath = "metrics.json"
	c.Save.LoadFile = false

	flag.StringVar(&c.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&c.Save.Interval, "i", 300, "interval save to file")
	flag.StringVar(&c.Save.FilePath, "f", "metrics.json", "file name save to metrics")
	flag.BoolVar(&c.Save.LoadFile, "r", false, "load old metrics")
	flag.Parse()

	if v, ok := os.LookupEnv("ADDRESS"); ok {
		c.FlagRunAddr = v
	}
	if v, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		if b, err := strconv.Atoi(v); err == nil {
			c.Save.Interval = b
		} else {
			logger.Infof("невалидный STORE_INTERVAL=%q: %v", v, err)
		}
	}
	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		c.Save.FilePath = v
	}
	if v, ok := os.LookupEnv("RESTORE"); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			c.Save.LoadFile = b
		} else {
			logger.Infof("невалидный STORE_INTERVAL=%q: %v", v, err)
		}
	}
	return c
}
