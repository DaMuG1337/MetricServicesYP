package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp_checkHandler(t *testing.T) {
	app := &App{storage: &MemStorage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	},
	}

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
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, tt.path, nil)
			app.checkHandler(w, r)
			if w.Code != tt.wantCode {
				t.Errorf("Не так = %v, Надо бы так %v", w.Code, tt.wantCode)
			}
		})
	}
}
