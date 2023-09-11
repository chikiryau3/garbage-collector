package memstorage

type StorageData map[string]any

type MemStorage interface {
	WriteMetric(name string, value any) error
	ReadMetric(name string) (any, bool)
	GetData() (*StorageData, error)
}

type storage struct {
	data StorageData
}

func (s *storage) WriteMetric(name string, value any) error {
	s.data[name] = value

	return nil
}

func (s *storage) ReadMetric(name string) (any, bool) {
	value, ok := s.data[name]

	return value, ok
}

func (s *storage) GetData() (*StorageData, error) {
	return &s.data, nil
}

func New() MemStorage {
	return &storage{
		data: map[string]any{},
	}
}
