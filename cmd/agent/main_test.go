package main

import (
	"reflect"
	"runtime"
	"testing"
)

func Test_addMetric(t *testing.T) {
	tests := []struct {
		name        string
		met         runtime.MemStats
		m           *Metrics
		wantMetrics *Metrics
	}{
		{
			name: "Тест 1 проверка на добавление метрик",
			met:  runtime.MemStats{Alloc: 100},
			m: &Metrics{
				gauge:   map[string]float64{},
				counter: map[string]int64{},
			},
			wantMetrics: &Metrics{
				gauge: map[string]float64{
					"Alloc":         100,
					"BuckHashSys":   0,
					"Frees":         0,
					"GCCPUFraction": 0,
					"GCSys":         0,
					"HeapAlloc":     0,
					"HeapIdle":      0,
					"HeapInuse":     0,
					"HeapObjects":   0,
					"HeapReleased":  0,
					"HeapSys":       0,
					"LastGC":        0,
					"Lookups":       0,
					"MCacheInuse":   0,
					"MCacheSys":     0,
					"MSpanInuse":    0,
					"MSpanSys":      0,
					"Mallocs":       0,
					"NextGC":        0,
					"NumForcedGC":   0,
					"NumGC":         0,
					"OtherSys":      0,
					"PauseTotalNs":  0,
					"StackInuse":    0,
					"StackSys":      0,
					"Sys":           0,
					"TotalAlloc":    0,
				},
				counter: map[string]int64{"PollCount": 1},
			},
		},
		{
			name: "Тест 2 Gauge заменяется, а не накапливается",
			met:  runtime.MemStats{Alloc: 999},
			m: &Metrics{
				gauge:   map[string]float64{"Alloc": 100},
				counter: map[string]int64{},
			},
			wantMetrics: &Metrics{
				gauge: map[string]float64{
					"Alloc":         999,
					"BuckHashSys":   0,
					"Frees":         0,
					"GCCPUFraction": 0,
					"GCSys":         0,
					"HeapAlloc":     0,
					"HeapIdle":      0,
					"HeapInuse":     0,
					"HeapObjects":   0,
					"HeapReleased":  0,
					"HeapSys":       0,
					"LastGC":        0,
					"Lookups":       0,
					"MCacheInuse":   0,
					"MCacheSys":     0,
					"MSpanInuse":    0,
					"MSpanSys":      0,
					"Mallocs":       0,
					"NextGC":        0,
					"NumForcedGC":   0,
					"NumGC":         0,
					"OtherSys":      0,
					"PauseTotalNs":  0,
					"StackInuse":    0,
					"StackSys":      0,
					"Sys":           0,
					"TotalAlloc":    0,
				},
				counter: map[string]int64{"PollCount": 1},
			},
		},
		{
			name: "Тест 3 PollCount увеличивается при повторном вызове",
			met:  runtime.MemStats{},
			m: &Metrics{
				gauge:   map[string]float64{},
				counter: map[string]int64{"PollCount": 5},
			},
			wantMetrics: &Metrics{
				gauge: map[string]float64{
					"Alloc":         0,
					"BuckHashSys":   0,
					"Frees":         0,
					"GCCPUFraction": 0,
					"GCSys":         0,
					"HeapAlloc":     0,
					"HeapIdle":      0,
					"HeapInuse":     0,
					"HeapObjects":   0,
					"HeapReleased":  0,
					"HeapSys":       0,
					"LastGC":        0,
					"Lookups":       0,
					"MCacheInuse":   0,
					"MCacheSys":     0,
					"MSpanInuse":    0,
					"MSpanSys":      0,
					"Mallocs":       0,
					"NextGC":        0,
					"NumForcedGC":   0,
					"NumGC":         0,
					"OtherSys":      0,
					"PauseTotalNs":  0,
					"StackInuse":    0,
					"StackSys":      0,
					"Sys":           0,
					"TotalAlloc":    0,
				},
				counter: map[string]int64{"PollCount": 6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addMetric(tt.met, tt.m)
			delete(tt.m.gauge, "RandomValue")
			if !reflect.DeepEqual(tt.m, tt.wantMetrics) {
				t.Errorf("Ожидалось: %v, получено: %v", tt.wantMetrics, tt.m)
			}
		})
	}
}
