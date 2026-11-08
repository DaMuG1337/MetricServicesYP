package filesave

import (
	"encoding/json"
	"os"
	"time"

	"github.com/DaMuG1337/MetricServices/internal/handler"
	"github.com/DaMuG1337/MetricServices/internal/logger"
	"github.com/DaMuG1337/MetricServices/internal/models"
)

func SaveToFile(filename string, metrics []models.Metrics) error {
	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0666)
}

func LoadToFile(filename string) ([]models.Metrics, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var metrics []models.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return nil, err
	}
	return metrics, nil
}

func RunSaveMetrics(filePath string, interval int, storage handler.Storage) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics := storage.GetAllMetrics()
		if err := SaveToFile(filePath, metrics); err != nil {
			logger.Log.Errorw("Ошибка сохранения", "ошибочка:", err)
		}
	}
}

func LoadFileMetrics(filename string, storage handler.Storage) error {
	metrics, err := LoadToFile(filename)
	if err != nil {
		return err
	}
	for _, m := range metrics {
		switch m.MType {
		case models.Gauge:
			if m.Value != nil {
				storage.SetGauge(m.ID, *m.Value)
			}
		case models.Counter:
			if m.Delta != nil {
				storage.AddCounter(m.ID, *m.Delta)
			}
		}
	}
	return nil
}

func SyncSaveFile(filePath string, interval int, storage handler.Storage) func() {
	if interval != 0 {
		return nil
	}
	return func() {
		if err := SaveToFile(filePath, storage.GetAllMetrics()); err != nil {
			logger.Log.Errorw("Ошибка синхронного сохранения", err)
		}
	}
}
