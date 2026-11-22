package db

import (
	"database/sql"

	"go.uber.org/zap"
)

func NewBD(dsn string, logger *zap.SugaredLogger) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error("Ошибка открытия к БД", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.Error("Ошибка подключения к БД", err)
		db.Close()
		return nil, err
	}

	return db, nil
}
