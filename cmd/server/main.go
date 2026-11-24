package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/DaMuG1337/MetricServices/internal/config"
	db "github.com/DaMuG1337/MetricServices/internal/database"
	"github.com/DaMuG1337/MetricServices/internal/database/migrations"
	filesave "github.com/DaMuG1337/MetricServices/internal/file_save"
	"github.com/DaMuG1337/MetricServices/internal/gzip"
	"github.com/DaMuG1337/MetricServices/internal/handler"
	logpkg "github.com/DaMuG1337/MetricServices/internal/logger"
	"github.com/DaMuG1337/MetricServices/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	logger, err := logpkg.New()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.ServerFlags(logger)

	var metricsStorage handler.Storage
	var database *sql.DB
	if cfg.DSN != "" {
		database, err = db.NewBD(cfg.DSN, logger)
		if err != nil {
			logger.Errorw("не удалось подключиться к БД", "error", err)
		} else {
			defer database.Close()
			if err := migrations.Migrate(database, logger); err != nil {
				logger.Errorw("ошибка применения миграций", "error", err)
				database.Close()
				database = nil
			} else {
				metricsStorage = storage.NewDBStorage(database, logger)
			}
		}
	}
	if metricsStorage == nil {
		metricsStorage, err = storage.NewMemStorage(cfg.Save.FilePath, cfg.Save.LoadFile)
		if err != nil {
			logger.Errorw("Ошибка конструктора метрик", "error", err)
		} else if cfg.Save.LoadFile {
			logger.Infow("Метрики загружены из файла", "path", cfg.Save.FilePath)
		}
	}

	if err := run(cfg, metricsStorage, database, logger); err != nil {
		logger.Errorw("Ошибка запуска сервера", "error", err)
	}
}

func run(cfg config.ConfigServer, metricsStorage handler.Storage, database *sql.DB, logger *zap.SugaredLogger) error {
	var onUpdate func()
	if database == nil {
		onUpdate = filesave.SyncSaveFile(cfg.Save.FilePath, cfg.Save.Interval, metricsStorage, logger)
	}
	if database == nil && cfg.Save.Interval > 0 {
		go filesave.RunSaveMetrics(cfg.Save.FilePath, cfg.Save.Interval, metricsStorage, logger)
	}
	app := handler.NewApp(metricsStorage, onUpdate)
	r := chi.NewRouter()

	r.Use(logpkg.LogResponse(logger))
	r.Use(gzip.GzipMiddleware)

	r.Post("/update/{type}/{name}/{value}", app.PostMetrics)
	r.Post("/update", app.PostMetricsJSON)
	r.Post("/update/", app.PostMetricsJSON)
	r.Post("/value", app.GetMetricsJSON)
	r.Post("/value/", app.GetMetricsJSON)
	r.Get("/value/{type}/{name}", app.GetMetrics)
	r.Get("/", app.GetAllMetrics)
	r.Get("/ping", app.CheckDB(database))

	return http.ListenAndServe(cfg.FlagRunAddr, r)
}
