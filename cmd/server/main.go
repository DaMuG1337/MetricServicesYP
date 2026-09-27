package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.checkHandler)
	return http.ListenAndServe(":8080", mux)
}

func (a *App) checkHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != MethodPost {
		http.Error(w, "Щас пока только POST", http.StatusMethodNotAllowed)
		return
	}
	partsMetric := strings.Split(r.URL.Path, "/")
	if len(partsMetric) != 5 || partsMetric[0] != "" || partsMetric[1] != "update" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	typeMetric := partsMetric[2]
	nameMetric := partsMetric[3]
	valueMetric := partsMetric[4]

	if nameMetric == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch typeMetric {
	case "gauge":
		val, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			fmt.Println("Ошибка преобразования:", err)
			return
		}
		a.storage.gauge[nameMetric] = val
		w.WriteHeader(http.StatusOK)
		fmt.Println(a.storage)
		return
	case "counter":
		val, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			fmt.Println("Ошибка преобразования:", err)
			return
		}
		a.storage.counter[nameMetric] += val
		w.WriteHeader(http.StatusOK)
		fmt.Println(a.storage)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Printf("Path: %q\n", r.URL.Path)
	w.WriteHeader(http.StatusOK)
}
