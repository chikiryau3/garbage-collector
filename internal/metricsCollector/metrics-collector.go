package metricscollector

// MetricsCollector интерфейс, содержащий бизнес-логику нашего сервиса
// это пока выглядит как прокси к стораджу (но все же логика формирования данных к самому стораджу отношения не имеет)
// когда будет БД, станет понятно зачем эта штука
//
// также делал все так, будто у нас может быть много разных метрик
// для каждой своя логика записи -- эту логику лучше деражть в отдельном от стораджа слое
type MetricsCollector interface {
	SetGauge(name string, value float64) (float64, error)
	SetCount(name string, value int64) (int64, error)

	ReadStorage() (*StorageData, error)
	GetMetric(mtype string, name string) (any, error)

	SetBatch(batch []Metrics) (*[]Metrics, error)
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
type StorageData map[string]any

// UPD: перенес объявление интерфейса в коллектор, тк он общий для разных storage

type Storage interface {
	WriteMetric(mtype string, name string, value any) error
	ReadMetric(mtype string, name string) (any, bool)
	GetData() (*StorageData, error)
}

type metricsCollector struct {
	storage Storage
}

func New(s Storage) MetricsCollector {
	return &metricsCollector{
		storage: s,
	}
}
