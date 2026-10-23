package main

import (
	"log"
	"net/http"

	"github.com/DaMuG1337/MetricServices/internal/config"
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
	app := handler.NewApp(storage)
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", logger.LogResponse(http.HandlerFunc(app.PostMetrics)))
	r.Get("/value/{type}/{name}", logger.LogResponse(http.HandlerFunc(app.GetMetrics)))
	r.Get("/", logger.LogResponse(http.HandlerFunc(app.GetAllMetrics)))
	return http.ListenAndServe(cfg.FlagRunAddr, r)
}
