package storage

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/DaMuG1337/MetricServices/internal/models"
)

type MemStorage struct {
	mu      sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage(filePath string, restore bool) (*MemStorage, error) {
	m := &MemStorage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	}
	if restore && filePath != "" {
		if err := m.LoadFromFile(filePath); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *MemStorage) SetGauge(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauge[name] = value
}

func (m *MemStorage) GetGauge(name string) (float64, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.gauge[name]
	return val, ok
}

func (m *MemStorage) GetCounter(name string) (int64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.counter[name]
	return v, ok
}

func (m *MemStorage) SetCounter(name string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter[name] = value
}

func (m *MemStorage) AddCounter(name string, val int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter[name] += val
}

func (m *MemStorage) GetAllGauges() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]float64, len(m.gauge))
	for k, v := range m.gauge {
		result[k] = v
	}
	return result
}

func (m *MemStorage) GetAllCounters() map[string]int64 {
	result := make(map[string]int64, len(m.counter))
	for k, v := range m.counter {
		result[k] = v
	}
	return result
}

func (m *MemStorage) GetAllMetrics() []models.Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []models.Metrics
	for name, value := range m.gauge {
		v := value
		result = append(result, models.Metrics{
			ID:    name,
			MType: models.Gauge,
			Value: &v,
		})
	}
	for name, value := range m.counter {
		v := value
		result = append(result, models.Metrics{
			ID:    name,
			MType: models.Counter,
			Delta: &v,
		})
	}
	return result
}

func (m *MemStorage) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var metrics []models.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, metric := range metrics {
		switch metric.MType {
		case models.Gauge:
			if metric.Value != nil {
				m.gauge[metric.ID] = *metric.Value
			}
		case models.Counter:
			if metric.Delta != nil {
				m.counter[metric.ID] += *metric.Delta
			}
		}
	}
	return nil
}
