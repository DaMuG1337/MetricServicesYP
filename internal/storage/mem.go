package storage

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	}
}

func (m *MemStorage) SetGauge(name string, value float64) {
	m.gauge[name] = value
}

func (m *MemStorage) GetGauge(name string) (float64, bool) {
	val, ok := m.gauge[name]
	return val, ok
}

func (m *MemStorage) GetCounter(name string) (int64, bool) {
	v, ok := m.counter[name]
	return v, ok
}

func (m *MemStorage) SetCounter(name string, value int64) {
	m.counter[name] = value
}

func (m *MemStorage) AddCounter(name string, val int64) {
	m.counter[name] += val
}

func (m *MemStorage) GetAllGauges() map[string]float64 {
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
