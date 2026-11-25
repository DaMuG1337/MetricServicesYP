package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func NewBD(dsn string, logger *zap.SugaredLogger) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Errorw("Ошибка открытия к БД", "error", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.Errorw("Ошибка подключения к БД", "error", err)
		db.Close()
		return nil, err
	}

	return db, nil
}
