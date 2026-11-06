package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DaMuG1337/MetricServices/internal/handler"
	"github.com/DaMuG1337/MetricServices/internal/storage"
	"github.com/go-chi/chi/v5"
)

func TestApp_postMetrics(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		wantCode int
	}{
		{" Test 1: Проверочка gauge", http.MethodPost, "/update/gauge/testGauge/123.45", http.StatusOK},
		{" Test 2: Проверочка counter", http.MethodPost, "/update/counter/testCounter/10", http.StatusOK},
		{" Test 3: Мисс метод", http.MethodGet, "/update/gauge/testGauge/123.45", http.StatusMethodNotAllowed},
		{" Test 4: Неверный формат пути", http.MethodPost, "/invalid/path/format", http.StatusNotFound},
		{" Test 5: Пустое имя метрики", http.MethodPost, "/update/gauge//123.45", http.StatusNotFound},
		{" Test 6: Неверное значение gauge", http.MethodPost, "/update/gauge/testGauge/invalidValue", http.StatusBadRequest},
		{" Test 7: Неверное значение counter", http.MethodPost, "/update/counter/testCounter/invalidValue", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := storage.NewMemStorage()
			app := handler.NewApp(storage)

			r := chi.NewRouter()
			r.Post("/update/{type}/{name}/{value}", app.PostMetrics)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			r.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("Не так = %v, Надо бы так %v", w.Code, tt.wantCode)
				return
			}
		})
	}
}

func TestApp_getMetrics(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		path      string
		wantCode  int
		wantValue string
	}{
		{"Первая проверка", http.MethodGet, "/value/gauge/Alloc", http.StatusOK, "123"},
		{"Вторая проверка", http.MethodGet, "/value/counter/PollCount", http.StatusOK, "42"},
		{"Третья попытка", http.MethodPost, "/value/gauge/Alloc", http.StatusMethodNotAllowed, ""},
		{"Четвёртая проверка", http.MethodGet, "/value/gauge/bobo", http.StatusNotFound, "Метрика по имени не найдена\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := storage.NewMemStorage()
			storage.SetGauge("Alloc", 123)
			storage.SetCounter("PollCount", 42)
			app := handler.NewApp(storage)

			r := chi.NewRouter()
			r.Get("/value/{type}/{name}", app.GetMetrics)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			r.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("Ожидался статус %d, получен %d", tt.wantCode, w.Code)
				return
			}
			if w.Body.String() != tt.wantValue {
				t.Errorf("Ожидалось тело %q, получено %q", tt.wantValue, w.Body.String())
				return
			}
		})
	}
}

func TestApp_postJSONMetrics(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		expectedCode int
	}{
		{"Первая проверка", http.MethodPost, `{"id": "Alloc", "type": "gauge", "value": 123}`, http.StatusOK},
		{"valid counter", http.MethodPost, `{"id":"PollCount","type":"counter","delta":5}`, http.StatusOK},
		{"invalid json", http.MethodPost, `{invalid`, http.StatusBadRequest},
		{"empty id", http.MethodPost, `{"id":"","type":"gauge","value":123}`, http.StatusNotFound},
		{"nil value", http.MethodPost, `{"id":"Alloc","type":"gauge"}`, http.StatusBadRequest},
		{"nil delta", http.MethodPost, `{"id":"PollCount","type":"counter"}`, http.StatusBadRequest},
		{"unknown type", http.MethodPost, `{"id":"Alloc","type":"unknown","value":123}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := storage.NewMemStorage()
			app := handler.NewApp(storage)
			r := chi.NewRouter()
			r.Post("/update", app.PostMetricsJSON)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "/update", strings.NewReader(tt.body))
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("code = %d, want %d", w.Code, tt.expectedCode)
			}
		})
	}
}

func TestApp_getJSONMetrics(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		setup        func(s handler.Storage)
		expectedCode int
		expectedBody string
	}{
		{
			name:   "valid gauge",
			method: http.MethodPost,
			body:   `{"id":"Alloc","type":"gauge"}`,
			setup: func(s handler.Storage) {
				s.SetGauge("Alloc", 123.45)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"Alloc","type":"gauge","value":123.45}`,
		},
		{
			name:   "valid counter",
			method: http.MethodPost,
			body:   `{"id":"PollCount","type":"counter"}`,
			setup: func(s handler.Storage) {
				s.AddCounter("PollCount", 5)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"PollCount","type":"counter","delta":5}`,
		},
		{
			name:         "metric not found",
			method:       http.MethodPost,
			body:         `{"id":"Unknown","type":"gauge"}`,
			setup:        func(s handler.Storage) {},
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
		{
			name:         "invalid json",
			method:       http.MethodPost,
			body:         `{invalid`,
			setup:        func(s handler.Storage) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "empty id",
			method:       http.MethodPost,
			body:         `{"id":"","type":"gauge"}`,
			setup:        func(s handler.Storage) {},
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
		{
			name:         "unknown type",
			method:       http.MethodPost,
			body:         `{"id":"Alloc","type":"unknown"}`,
			setup:        func(s handler.Storage) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := storage.NewMemStorage()
			tt.setup(storage)

			app := handler.NewApp(storage)

			r := chi.NewRouter()
			r.Post("/value", app.GetMetricsJSON)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "/value", strings.NewReader(tt.body))
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("code = %d, want %d", w.Code, tt.expectedCode)
			}
			if tt.expectedBody != "" && strings.TrimSpace(w.Body.String()) != tt.expectedBody {
				t.Errorf("body = %q, want %q", w.Body.String(), tt.expectedBody)
			}
		})
	}
}
