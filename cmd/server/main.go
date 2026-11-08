package main

import (
	"log"
	"net/http"

	"github.com/DaMuG1337/MetricServices/internal/config"
	filesave "github.com/DaMuG1337/MetricServices/internal/file_save"
	"github.com/DaMuG1337/MetricServices/internal/gzip"
	"github.com/DaMuG1337/MetricServices/internal/handler"
	"github.com/DaMuG1337/MetricServices/internal/logger"
	"github.com/DaMuG1337/MetricServices/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := logger.Initialize(); err != nil {
		panic(err)
	}
	cfg := config.ServerFlags()
	storage := storage.NewMemStorage()
	if err := run(cfg, storage); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.ConfigServer, storage handler.Storage) error {
	onUpdate := filesave.SyncSaveFile(cfg.Save.FilePath, cfg.Save.Interval, storage)
	app := handler.NewApp(storage, onUpdate)
	r := chi.NewRouter()
	if cfg.Save.LoadFile {
		if err := filesave.LoadFileMetrics(cfg.Save.FilePath, storage); err != nil {
			logger.Log.Errorw("Ошибка загрузки метрик из файла", err)
		} else {
			logger.Log.Infow("Метрики загружены из файла", cfg.Save.FilePath)
		}
	}
	r.Post("/update/{type}/{name}/{value}", logger.LogResponse(gzip.GzipMiddleware(http.HandlerFunc(app.PostMetrics))))
	r.Post("/update", logger.LogResponse(gzip.GzipMiddleware(http.HandlerFunc(app.PostMetricsJSON))))
	r.Post("/update/", logger.LogResponse(gzip.GzipMiddleware(http.HandlerFunc(app.PostMetricsJSON))))
	r.Post("/value", logger.LogResponse(gzip.GzipMiddleware(http.HandlerFunc(app.GetMetricsJSON))))
	r.Post("/value/", logger.LogResponse(gzip.GzipMiddleware(http.HandlerFunc(app.GetMetricsJSON))))
	r.Get("/value/{type}/{name}", logger.LogResponse(gzip.GzipMiddleware(http.HandlerFunc(app.GetMetrics))))
	r.Get("/", logger.LogResponse(gzip.GzipMiddleware(http.HandlerFunc(app.GetAllMetrics))))

	if cfg.Save.Interval > 0 {
		go filesave.RunSaveMetrics(cfg.Save.FilePath, cfg.Save.Interval, storage)
	}
	return http.ListenAndServe(cfg.FlagRunAddr, r)
}
