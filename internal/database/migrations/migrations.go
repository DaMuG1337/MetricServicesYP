package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"go.uber.org/zap"
)

//go:embed *.up.sql
var migrationFiles embed.FS

func Migrate(database *sql.DB, logger *zap.SugaredLogger) error {
	if database == nil {
		logger.Error("База данных пустая")
		return fmt.Errorf("database is nil")
	}

	dir, err := migrationFiles.ReadDir(".")
	if err != nil {
		logger.Error("Ошибка чтения файлов миграции")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Ошибка BeginTx")
		return err
	}
	defer tx.Rollback()

	for _, entry := range dir {
		file, err := migrationFiles.ReadFile(entry.Name())
		if err != nil {
			logger.Error("Ошибка чтения файла")
			return err
		}
		if _, err := tx.ExecContext(ctx, string(file)); err != nil {
			logger.Error("Ошибка выполнения миграции")
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		logger.Error("Ошибка выполнения Commit")
		return err
	}
	return nil
}
