package filesave

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/DaMuG1337/MetricServices/internal/handler"
	"github.com/DaMuG1337/MetricServices/internal/models"
	"go.uber.org/zap"
)

func SaveToFile(filename string, metrics []models.Metrics) error {
	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return err
	}
	// временный файл
	dir := filepath.Dir(filename)
	tmp, err := os.CreateTemp(dir, "metrics-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	// если что-то упадёт — удаляем временный
	defer func() {
		if err != nil {
			os.Remove(tmpName)
		}
	}()
	// пишем
	if _, err = tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	// синхронизируем на диск
	if err = tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	// закрываем
	if err = tmp.Close(); err != nil {
		return err
	}
	// атомарно заменяем
	if err = os.Rename(tmpName, filename); err != nil {
		return err
	}
	return nil
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

func RunSaveMetrics(filePath string, interval int, storage handler.Storage, logger *zap.SugaredLogger) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics := storage.GetAllMetrics()
		if err := SaveToFile(filePath, metrics); err != nil {
			logger.Errorw("Ошибка сохранения", "ошибочка:", err)
		}
	}
}

func SyncSaveFile(filePath string, interval int, storage handler.Storage, logger *zap.SugaredLogger) func() {
	if interval != 0 {
		return nil
	}
	return func() {
		if err := SaveToFile(filePath, storage.GetAllMetrics()); err != nil {
			logger.Errorw("Ошибка синхронного сохранения", err)
		}
	}
}
