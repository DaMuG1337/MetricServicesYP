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
		{"Test 1: Проверка добавления метрик", runtime.MemStats{}, &Metrics{
			gauge:   map[string]float64{},
			counter: map[string]int64{},
		},
			&Metrics{
				gauge:   map[string]float64{"Alloc": 0.0, "BuckHashSys": 0.0, "Frees": 0.0, "GCCPUFraction": 0.0, "GCSys": 0.0, "HeapAlloc": 0.0, "HeapIdle": 0.0, "HeapInuse": 0.0, "HeapObjects": 0.0, "HeapReleased": 0.0, "HeapSys": 0.0, "LastGC": 0.0, "Lookups": 0.0, "MCacheInuse": 0.0, "MCacheSys": 0.0, "MSpanInuse": 0.0, "MSpanSys": 0.0, "Mallocs": 0.0},
				counter: map[string]int64{"PollCount": 1},
			},
		},
		{"Test 2: Проверка добавления метрик с уже существующими значениями", runtime.MemStats{}, &Metrics{
			gauge:   map[string]float64{"Alloc": 100.0},
			counter: map[string]int64{"PollCount": 5},
		}, &Metrics{
			gauge:   map[string]float64{"Alloc": 100.0},
			counter: map[string]int64{"PollCount": 6},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.m, tt.wantMetrics) {
				t.Errorf("Ожидалось: %v, получено: %v", tt.wantMetrics, tt.m)
			}
		})
	}
}
