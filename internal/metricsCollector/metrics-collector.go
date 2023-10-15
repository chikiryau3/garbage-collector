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
	GetMetric(name string) (any, error)
}

type StorageData map[string]any

// UPD: перенес объявление интерфейса в коллектор, тк он общий для разных storage

type Storage interface {
	WriteMetric(name string, value any) error
	ReadMetric(name string) (any, bool)
	GetData() (*StorageData, error)

	RunStorageDumper() <-chan error
	RestoreFromDump() error
}

type metricsCollector struct {
	storage Storage
}

func New(s Storage) MetricsCollector {
	return &metricsCollector{
		storage: s,
	}
}
