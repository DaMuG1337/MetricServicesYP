package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DaMuG1337/MetricServices/internal/models"
	"github.com/go-chi/chi/v5"
)

type App struct {
	storage Storage
}

func NewApp(storage Storage) *App {
	return &App{storage: storage}
}

func (a *App) PostMetricsJSON(w http.ResponseWriter, r *http.Request) {
	var req models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch req.MType {
	case models.Gauge:
		if req.Value == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		a.storage.SetGauge(req.ID, *req.Value)
		w.WriteHeader(http.StatusOK)
		return
	case models.Counter:
		if req.Delta == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		a.storage.AddCounter(req.ID, *req.Delta)
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (a *App) GetMetricsJSON(w http.ResponseWriter, r *http.Request) {
	var req models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.MType == "" || req.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch req.MType {
	case models.Gauge:
		val, ok := a.storage.GetGauge(req.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		}
		req.Value = &val
	case models.Counter:
		val, ok := a.storage.GetCounter(req.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		}
		req.Delta = &val
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}

func (a *App) PostMetrics(w http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "type")
	nameMetric := chi.URLParam(r, "name")
	if nameMetric == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	valueMetric := chi.URLParam(r, "value")

	switch typeMetric {
	case "gauge":
		val, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		a.storage.SetGauge(nameMetric, val)
		w.WriteHeader(http.StatusOK)
		return
	case "counter":
		val, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		a.storage.AddCounter(nameMetric, val)
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (a *App) GetMetrics(w http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "type")
	nameMetric := chi.URLParam(r, "name")

	switch typeMetric {
	case "gauge":
		val, ok := a.storage.GetGauge(nameMetric)
		if !ok {
			http.Error(w, "Метрика по имени не найдена", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "%v", val)

	case "counter":
		val, ok := a.storage.GetCounter(nameMetric)
		if !ok {
			http.Error(w, "Метрика по имени не найдена", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "%v", val)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (a *App) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<html><body>"))
	w.Write([]byte("<h1>Метрики</h1>"))
	gauges := a.storage.GetAllGauges()
	counters := a.storage.GetAllCounters()
	w.Write([]byte("<h2>Gauge</h2>"))
	for name, value := range gauges {
		w.Write([]byte(fmt.Sprintf("<p>%s : %v</p>", name, value)))
	}
	w.Write([]byte("<h2>Counter</h2>"))
	for name, value := range counters {
		w.Write([]byte(fmt.Sprintf("<p>%s : %v</p>", name, value)))
	}
	w.Write([]byte("</body></html>"))
}
