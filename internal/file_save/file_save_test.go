package filesave

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DaMuG1337/MetricServices/internal/handler"
	"github.com/DaMuG1337/MetricServices/internal/models"
	"github.com/DaMuG1337/MetricServices/internal/storage"
	"go.uber.org/zap"
)

func floatPtr(f float64) *float64 { return &f }
func int64Ptr(i int64) *int64     { return &i }

func TestSaveToFile(t *testing.T) {
	tests := []struct {
		name    string
		metrics []models.Metrics
		wantErr bool
	}{
		{
			name: "одна gauge-метрика",
			metrics: []models.Metrics{
				{ID: "Alloc", MType: models.Gauge, Value: floatPtr(123.45)},
			},
			wantErr: false,
		},
		{
			name: "одна counter-метрика",
			metrics: []models.Metrics{
				{ID: "PollCount", MType: models.Counter, Delta: int64Ptr(5)},
			},
			wantErr: false,
		},
		{
			name: "несколько метрик",
			metrics: []models.Metrics{
				{ID: "Alloc", MType: models.Gauge, Value: floatPtr(1.5)},
				{ID: "PollCount", MType: models.Counter, Delta: int64Ptr(10)},
			},
			wantErr: false,
		},
		{
			name:    "пустой список",
			metrics: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "metrics.json")

			err := SaveToFile(path, tt.metrics)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToFile() ошибка = %v, ожидалась = %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if _, err := os.Stat(path); err != nil {
					t.Errorf("Файл не создан: %v", err)
				}
			}
		})
	}
}

func TestLoadToFile(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(path string) error
		want    []models.Metrics
		wantErr bool
	}{
		{
			name: "файл с gauge-метрикой",
			setup: func(path string) error {
				return SaveToFile(path, []models.Metrics{
					{ID: "Alloc", MType: models.Gauge, Value: floatPtr(123.45)},
				})
			},
			want: []models.Metrics{
				{ID: "Alloc", MType: models.Gauge, Value: floatPtr(123.45)},
			},
			wantErr: false,
		},
		{
			name: "файл с counter-метрикой",
			setup: func(path string) error {
				return SaveToFile(path, []models.Metrics{
					{ID: "PollCount", MType: models.Counter, Delta: int64Ptr(5)},
				})
			},
			want: []models.Metrics{
				{ID: "PollCount", MType: models.Counter, Delta: int64Ptr(5)},
			},
			wantErr: false,
		},
		{
			name:    "файл не существует",
			setup:   func(path string) error { return nil },
			want:    nil,
			wantErr: true,
		},
		{
			name: "невалидный JSON",
			setup: func(path string) error {
				return os.WriteFile(path, []byte("{invalid"), 0666)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "metrics.json")

			if err := tt.setup(path); err != nil {
				t.Fatalf("Ошибка подготовки: %v", err)
			}

			got, err := LoadToFile(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadToFile() ошибка = %v, ожидалась = %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Fatalf("Количество метрик = %d, ожидалось %d", len(got), len(tt.want))
				}
				for i := range tt.want {
					if got[i].ID != tt.want[i].ID || got[i].MType != tt.want[i].MType {
						t.Errorf("Метрика = %+v, ожидалась %+v", got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestSyncSaveFile(t *testing.T) {
	tests := []struct {
		name     string
		interval int
		wantNil  bool
	}{
		{
			name:     "интервал 0 — возвращает функцию",
			interval: 0,
			wantNil:  false,
		},
		{
			name:     "интервал 5 — возвращает nil",
			interval: 5,
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st, err := storage.NewMemStorage("", false)
			if err != nil {
				t.Fatal(err)
			}
			log := zap.NewNop().Sugar()
			fn := SyncSaveFile("test.json", tt.interval, st, log)

			if (fn == nil) != tt.wantNil {
				t.Errorf("SyncSaveFile() = %T, ожидался nil = %v", fn, tt.wantNil)
			}
		})
	}
}

func TestSyncSaveFile_SavesFile(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(st handler.Storage)
		wantErr bool
	}{
		{
			name: "с gauge-метрикой",
			setup: func(st handler.Storage) {
				st.SetGauge("Alloc", 1.5)
			},
			wantErr: false,
		},
		{
			name: "с counter-метрикой",
			setup: func(st handler.Storage) {
				st.AddCounter("PollCount", 5)
			},
			wantErr: false,
		},
		{
			name:    "пустое хранилище",
			setup:   func(st handler.Storage) {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "metrics.json")
			st, err := storage.NewMemStorage("", false)
			if err != nil {
				t.Fatal(err)
			}
			tt.setup(st)

			log := zap.NewNop().Sugar()
			fn := SyncSaveFile(path, 0, st, log)
			if fn == nil {
				t.Fatal("SyncSaveFile() вернул nil, ожидалась функция")
			}
			fn()

			if _, err := os.Stat(path); err != nil {
				t.Errorf("Файл не создан: %v", err)
			}
		})
	}
}

func TestRunSaveMetrics(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропуск медленного теста")
	}

	tests := []struct {
		name     string
		interval int
		setup    func(st handler.Storage)
	}{
		{
			name:     "сохраняет gauge",
			interval: 1,
			setup: func(st handler.Storage) {
				st.SetGauge("Alloc", 1.5)
			},
		},
		{
			name:     "сохраняет counter",
			interval: 1,
			setup: func(st handler.Storage) {
				st.AddCounter("PollCount", 5)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "metrics.json")
			st, err := storage.NewMemStorage("", false)
			if err != nil {
				t.Fatal(err)
			}
			tt.setup(st)

			log := zap.NewNop().Sugar()
			go RunSaveMetrics(path, tt.interval, st, log)

			deadline := time.Now().Add(3 * time.Second)
			for time.Now().Before(deadline) {
				if _, err := os.Stat(path); err == nil {
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
			t.Error("Файл не создан за 3 секунды")
		})
	}
}
