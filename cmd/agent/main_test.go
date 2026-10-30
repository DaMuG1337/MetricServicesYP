package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DaMuG1337/MetricServices/internal/config"
	"github.com/DaMuG1337/MetricServices/internal/logger"
	"github.com/DaMuG1337/MetricServices/internal/models"
)

func floatPtr(f float64) *float64 {
	return &f
}

func int64Ptr(a int64) *int64 {
	return &a
}

func TestMakeJSONRequest_Success(t *testing.T) {
	if err := logger.Initialize(); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		metrics models.Metrics
	}{
		{
			name: "gauge positive",
			metrics: models.Metrics{
				ID:    "Alloc",
				MType: models.Gauge,
				Value: floatPtr(123.45),
			},
		},
		{
			name: "gauge zero",
			metrics: models.Metrics{
				ID:    "Alloc",
				MType: models.Gauge,
				Value: floatPtr(0),
			},
		},
		{
			name: "gauge negative",
			metrics: models.Metrics{
				ID:    "Alloc",
				MType: models.Gauge,
				Value: floatPtr(-42.5),
			},
		},
		{
			name: "counter positive",
			metrics: models.Metrics{
				ID:    "PollCount",
				MType: models.Counter,
				Delta: int64Ptr(5),
			},
		},
		{
			name: "counter zero",
			metrics: models.Metrics{
				ID:    "PollCount",
				MType: models.Counter,
				Delta: int64Ptr(0),
			},
		},
		{
			name: "counter large",
			metrics: models.Metrics{
				ID:    "PollCount",
				MType: models.Counter,
				Delta: int64Ptr(1_000_000),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedBody []byte
			var receivedContentType string
			var receivedMethod string

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				receivedContentType = r.Header.Get("Content-Type")
				receivedBody, _ = io.ReadAll(r.Body)
				w.WriteHeader(http.StatusOK)
			}))
			defer srv.Close()

			addr := strings.TrimPrefix(srv.URL, "http://")
			cfg := config.ConfigAgent{
				Server: config.ServerConfig{FlagRunAddr: addr},
			}

			err := makeJSONRequest(tt.metrics, cfg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if receivedMethod != http.MethodPost {
				t.Errorf("method = %s, want POST", receivedMethod)
			}
			if receivedContentType != "application/json" {
				t.Errorf("content-type = %s, want application/json", receivedContentType)
			}

			var got models.Metrics
			if err := json.Unmarshal(receivedBody, &got); err != nil {
				t.Fatalf("invalid json: %v", err)
			}
			if got.ID != tt.metrics.ID || got.MType != tt.metrics.MType {
				t.Errorf("id/type = %s/%s, want %s/%s", got.ID, got.MType, tt.metrics.ID, tt.metrics.MType)
			}
			if tt.metrics.Value != nil && (got.Value == nil || *got.Value != *tt.metrics.Value) {
				t.Errorf("value = %v, want %v", got.Value, tt.metrics.Value)
			}
			if tt.metrics.Delta != nil && (got.Delta == nil || *got.Delta != *tt.metrics.Delta) {
				t.Errorf("delta = %v, want %v", got.Delta, tt.metrics.Delta)
			}
		})
	}
}
