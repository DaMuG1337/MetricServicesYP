package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

type App struct {
	storage *MemStorage
}

func main() {
	parseFlags()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	app := &App{storage: &MemStorage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	},
	}
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", app.postMetrics)
	r.Get("/value/{type}/{name}", app.getMetrics)
	r.Get("/", app.getAllMetrics)
	return http.ListenAndServe(flagRunAddr, r)
}

func (a *App) postMetrics(w http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "type")
	nameMetric := chi.URLParam(r, "name")
	valueMetric := chi.URLParam(r, "value")

	switch typeMetric {
	case "gauge":
		val, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		a.storage.gauge[nameMetric] = val
		w.WriteHeader(http.StatusOK)
		fmt.Println(a.storage)
		return
	case "counter":
		val, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		a.storage.counter[nameMetric] += val
		w.WriteHeader(http.StatusOK)
		fmt.Println(a.storage)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (a *App) getMetrics(w http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "type")
	nameMetric := chi.URLParam(r, "name")

	switch typeMetric {
	case "gauge":
		val, ok := a.storage.gauge[nameMetric]
		if !ok {
			http.Error(w, "Метрика по имени не найдена", http.StatusNotFound)
		}
		fmt.Fprintf(w, "%v", val)

	case "counter":
		val, ok := a.storage.counter[nameMetric]
		if !ok {
			http.Error(w, "Метрика по имени не найдена", http.StatusNotFound)
		}
		fmt.Fprintf(w, "%v", val)
	}
}

func (a *App) getAllMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<html><body>"))
	w.Write([]byte("<h1>Метрики</h1>"))
	for name, value := range a.storage.gauge {
		w.Write([]byte(fmt.Sprintf("<p>%s : %v</p>", name, value)))
	}
	for name, value := range a.storage.counter {
		w.Write([]byte(fmt.Sprintf("<p>%s : %v</p>", name, value)))
	}
	w.Write([]byte("</body></html>"))
}
