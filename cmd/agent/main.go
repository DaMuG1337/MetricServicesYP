package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"runtime"
	"time"

	"github.com/DaMuG1337/MetricServices/internal/config"
	"github.com/DaMuG1337/MetricServices/internal/gzip/compress"
	"github.com/DaMuG1337/MetricServices/internal/logger"
	"github.com/DaMuG1337/MetricServices/internal/models"
	"go.uber.org/zap"
)

// инициализирую метрику с 2 типами
type Metrics struct {
	gauge   map[string]float64
	counter map[string]int64
}

// основная функция агента
func main() {
	m := Metrics{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	}

	logger, err := logger.New()
	if err != nil {
		logger.Error("Ошибка запуска логгера", err)
	}
	cfg := config.AgentFlags(logger)

	var met runtime.MemStats
	pollInterval := cfg.Poll.FlagPollInterval
	reportInterval := cfg.Poll.FlagReportInterval
	pollCount := 0

	for {
		runtime.ReadMemStats(&met)
		addMetric(met, &m)
		pollCount++

		if pollCount%(reportInterval/pollInterval) == 0 {
			rangeMetrics(&m, cfg, logger)
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

func rangeMetrics(m *Metrics, cfg config.ConfigAgent, logger *zap.SugaredLogger) {
	for name, value := range m.gauge {
		v := value
		if err := makeJSONRequest(models.Metrics{
			ID:    name,
			MType: models.Gauge,
			Value: &v,
		}, cfg, logger); err != nil {
			logger.Error("запрос выполнен с кодом %")
		}
	}
	for name, value := range m.counter {
		v := value
		if err := makeJSONRequest(models.Metrics{
			ID:    name,
			MType: models.Counter,
			Delta: &v,
		}, cfg, logger); err != nil {
			logger.Error("запрос выполнен с кодом %")
		}
	}
	m.counter["PollCount"] = 0 // обнуляю счетчик после отправки
}

func makeJSONRequest(m models.Metrics, cfg config.ConfigAgent, logger *zap.SugaredLogger) error {
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}
	compressed, err := compress.Compress(body)
	if err != nil {
		return err
	}
	start := time.Now()
	url := fmt.Sprintf("http://%s/update", cfg.Server.FlagRunAddr)

	logger.Infow("request",
		"URL", url,
		"Метод", "POST",
		"Время выполнения запроса", time.Since(start),
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(compressed))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("запрос выполнен с кодом %d", resp.StatusCode)
	}

	logger.Infow("response",
		"Статус", resp.StatusCode,
		"Размер", len(body),
		"Вермя выполнения", time.Since(start),
	)
	return nil
}
