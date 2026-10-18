package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"runtime"
	"time"

	"github.com/DaMuG1337/MetricServices/internal/config"
)

type Metrics struct {
	gauge   map[string]float64
	counter map[string]int64
}

func main() {
	m := Metrics{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	}
	cfg := config.AgentFlags()

	var met runtime.MemStats
	pollInterval := cfg.Poll.FlagPollInterval
	reportInterval := cfg.Poll.FlagReportInterval
	pollCount := 0

	for {
		runtime.ReadMemStats(&met)
		addMetric(met, &m)
		pollCount++

		if pollCount%(reportInterval/pollInterval) == 0 {
			rangeMetrics(&m, cfg)
		}

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func addMetric(met runtime.MemStats, m *Metrics) {
	m.counter["PollCount"]++                // увеличиваю на 1 при обновлении метрик
	m.gauge["RandomValue"] = rand.Float64() // что то просто

	m.gauge["Alloc"] = float64(met.Alloc)
	m.gauge["BuckHashSys"] = float64(met.BuckHashSys)
	m.gauge["Frees"] = float64(met.Frees)
	m.gauge["GCCPUFraction"] = float64(met.GCCPUFraction)
	m.gauge["GCSys"] = float64(met.GCSys)
	m.gauge["HeapAlloc"] = float64(met.HeapAlloc)
	m.gauge["HeapIdle"] = float64(met.HeapIdle)
	m.gauge["HeapInuse"] = float64(met.HeapInuse)
	m.gauge["HeapObjects"] = float64(met.HeapObjects)
	m.gauge["HeapReleased"] = float64(met.HeapReleased)
	m.gauge["HeapSys"] = float64(met.HeapSys)
	m.gauge["LastGC"] = float64(met.LastGC)
	m.gauge["Lookups"] = float64(met.Lookups)
	m.gauge["MCacheInuse"] = float64(met.MCacheInuse)
	m.gauge["MCacheSys"] = float64(met.MCacheSys)
	m.gauge["MSpanInuse"] = float64(met.MSpanInuse)
	m.gauge["MSpanSys"] = float64(met.MSpanSys)
	m.gauge["Mallocs"] = float64(met.Mallocs)
	m.gauge["NextGC"] = float64(met.NextGC)
	m.gauge["NumForcedGC"] = float64(met.NumForcedGC)
	m.gauge["NumGC"] = float64(met.NumGC)
	m.gauge["OtherSys"] = float64(met.OtherSys)
	m.gauge["PauseTotalNs"] = float64(met.PauseTotalNs)
	m.gauge["StackInuse"] = float64(met.StackInuse)
	m.gauge["StackSys"] = float64(met.StackSys)
	m.gauge["Sys"] = float64(met.Sys)
	m.gauge["TotalAlloc"] = float64(met.TotalAlloc)
}

func rangeMetrics(m *Metrics, cfg config.ConfigAgent) {
	for name, value := range m.gauge {
		if err := makeURLRequest("gauge", name, value, cfg); err != nil {
			log.Printf("Ошибка формирования запроса для %s: %v", name, err)
		}
	}
	for name, value := range m.counter {
		if err := makeURLRequest("counter", name, value, cfg); err != nil {
			log.Printf("Ошибка формирования запроса для %s: %v", name, err)
		}
	}
	m.counter["PollCount"] = 0 // обнуляю счетчик после отправки
}

func makeURLRequest(typeMetric, name string, value interface{}, cfg config.ConfigAgent) error {
	s := fmt.Sprintf("http://%s/update/%s/%s/%v", cfg.Server.FlagRunAddr, typeMetric, name, value)
	log.Print(s)
	response, err := http.Post(s, "text/plain", nil)
	if err != nil {
		log.Println("Ошибка:", err)
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Println("Сервер вернул:", response.Status)
	}
	return nil
}
