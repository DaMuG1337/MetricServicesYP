package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DaMuG1337/MetricServices/internal/config"
	"github.com/DaMuG1337/MetricServices/internal/handler"
	"github.com/DaMuG1337/MetricServices/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.ServerFlags()
	storage := storage.NewMemStorage()
	if err := run(cfg, storage); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.ConfigServer, storage handler.Storage) error {
	app := handler.NewApp(storage)
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", app.PostMetrics)
	r.Get("/value/{type}/{name}", app.GetMetrics)
	r.Get("/", app.GetAllMetrics)
	fmt.Println(cfg.FlagRunAddr)
	return http.ListenAndServe(cfg.FlagRunAddr, r)
}
