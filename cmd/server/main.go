package main

import (
	"log"
	"net/http"

	"github.com/DaMuG1337/MetricServices/internal/config"
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
	storage, err := storage.NewMemStorage(cfg.Save.FilePath, cfg.Save.LoadFile)
	if err != nil {
		logger.Error("Ошибка конструктора метрик", err)
	} else if cfg.Save.LoadFile {
		logger.Infow("Метрики загружены из файла", cfg.Save.FilePath)
	}
	if err := run(cfg, storage, logger); err != nil {
		logger.Error("Ошибка запуска сервера", err)
	}
}

func run(cfg config.ConfigServer, storage handler.Storage, logger *zap.SugaredLogger) error {
	onUpdate := filesave.SyncSaveFile(cfg.Save.FilePath, cfg.Save.Interval, storage, logger)
	app := handler.NewApp(storage, onUpdate)
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

	if cfg.Save.Interval > 0 {
		go filesave.RunSaveMetrics(cfg.Save.FilePath, cfg.Save.Interval, storage, logger)
	}
	return http.ListenAndServe(cfg.FlagRunAddr, r)
}
