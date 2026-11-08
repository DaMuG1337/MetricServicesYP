package handler

import "github.com/DaMuG1337/MetricServices/internal/models"

type Storage interface {
	SetGauge(name string, value float64)
	GetGauge(name string) (float64, bool)
	GetAllGauges() map[string]float64

	SetCounter(name string, value int64)
	AddCounter(name string, value int64)
	GetCounter(name string) (int64, bool)
	GetAllCounters() map[string]int64
	GetAllMetrics() []models.Metrics
}
