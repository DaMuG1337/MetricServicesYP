package main

import (
	"net/http"
	"net/http/httptest"
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
	storage := storage.NewMemStorage()
	app := handler.NewApp(storage)

	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", app.GetMetrics)

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
		{"Первая проверка", http.MethodGet, "/value/gauge/bobo", http.StatusNotFound, "Метрика по имени не найдена\n0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
