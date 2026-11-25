package storage

import (
	"database/sql"

	"github.com/DaMuG1337/MetricServices/internal/models"
	"go.uber.org/zap"
)

type DBStorage struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewDBStorage(db *sql.DB, logger *zap.SugaredLogger) *DBStorage {
	return &DBStorage{
		db:     db,
		logger: logger,
	}
}

func (s *DBStorage) SetGauge(name string, value float64) {
	_, err := s.db.Exec(`
	INSERT INTO metrics (id, type, value)
	VALUES ($1, $2, $3)
	ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value
	`, name, models.Gauge, value)

	if err != nil {
		s.logger.Errorw("ошибка записи gauge", "error", err)
	}
}

func (s *DBStorage) GetGauge(name string) (float64, bool) {
	var value sql.NullFloat64

	err := s.db.QueryRow(`SELECT value FROM metrics WHERE id = $1 AND type = $2`, name, models.Gauge).Scan(&value)
	if err != nil {
		if err != sql.ErrNoRows {
			s.logger.Errorw("Ошибка получения Gauge", "error", err)
		}
		return 0, false
	}
	if !value.Valid {
		return 0, false
	}
	return value.Float64, true
}

func (s *DBStorage) GetAllGauges() map[string]float64 {
	result := make(map[string]float64)
	rows, err := s.db.Query(`SELECT id, value FROM metrics WHERE type = $1 AND value IS NOT NULL`, models.Gauge)
	if err != nil {
		s.logger.Errorw("Ошибка получения ВСЕХ Gauge", "error", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var value float64
		if err := rows.Scan(&name, &value); err != nil {
			s.logger.Errorw("Ошибка чтения всех Gauge", "error", err)
			continue
		}
		result[name] = value
	}
	if err := rows.Err(); err != nil {
		s.logger.Errorw("Ошибка при выполнении запроса", "error", err)
	}
	return result
}

func (s *DBStorage) SetCounter(name string, value int64) {
	_, err := s.db.Exec(`
	INSERT INTO metrics (id, type, delta)
	VALUES ($1, $2, $3)
	ON CONFLICT (id) DO UPDATE SET delta = EXCLUDED.delta
	`, name, models.Counter, value)

	if err != nil {
		s.logger.Errorw("ошибка записи counter", "error", err)
	}
}

func (s *DBStorage) GetCounter(name string) (int64, bool) {
	var value sql.NullInt64

	err := s.db.QueryRow(`SELECT delta FROM metrics WHERE id = $1 AND type = $2`, name, models.Counter).Scan(&value)
	if err != nil {
		if err != sql.ErrNoRows {
			s.logger.Errorw("Ошибка получения Counter", "error", err)
		}
		return 0, false
	}
	if !value.Valid {
		return 0, false
	}
	return value.Int64, true
}

func (s *DBStorage) AddCounter(name string, value int64) {
	_, err := s.db.Exec(`
		INSERT INTO metrics (id, type, delta)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET type = EXCLUDED.type, delta = metrics.delta + EXCLUDED.delta, value = NULL`,
		name, models.Counter, value)
	if err != nil {
		s.logger.Errorw("Ошибка при сложении Counter", "error", err)
	}
}

func (s *DBStorage) GetAllCounters() map[string]int64 {
	result := make(map[string]int64)
	rows, err := s.db.Query(`SELECT id, delta FROM metrics WHERE type = $1 AND delta IS NOT NULL`, models.Counter)
	if err != nil {
		s.logger.Errorw("Ошибка преобразования ВСЕХ Counter", "error", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var value int64
		if err := rows.Scan(&name, &value); err != nil {
			s.logger.Errorw("Ошибка преобразования ВСЕХ Counter", "error", err)
			continue
		}
		result[name] = value
	}
	if err := rows.Err(); err != nil {
		s.logger.Errorw("Ошибка при выполнении запроса", "error", err)
	}
	return result
}

func (s *DBStorage) GetAllMetrics() []models.Metrics {
	result := []models.Metrics{}

	gauges := s.GetAllGauges()
	for name, val := range gauges {
		v := val
		result = append(result, models.Metrics{
			ID:    name,
			MType: models.Gauge,
			Value: &v,
		})
	}

	counters := s.GetAllCounters()
	for name, val := range counters {
		v := val
		result = append(result, models.Metrics{
			ID:    name,
			MType: models.Counter,
			Delta: &v,
		})
	}
	return result
}
