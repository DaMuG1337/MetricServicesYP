package db

import (
	"flag"
	"os"

	"go.uber.org/zap"
)

type DBConfig struct {
	DSN string
}

func DBFlag(logger *zap.SugaredLogger) DBConfig {
	var db DBConfig

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	CommandLine.StringVar(&db.DSN, "d", "", "database connection string")
	if err := CommandLine.Parse(os.Args[1:]); err != nil {
		logger.Error("Ошибка парсинга флагов:", err)
	}
	if db.DSN == "" {
		if envDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
			db.DSN = envDSN
		}
	}

	return db
}
