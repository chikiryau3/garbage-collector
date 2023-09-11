package memStorage

type MemStorage interface {
	WriteMetric(name string, value any) error
	ReadMetric(name string) (any, bool)
}

type storage struct {
	data map[string]any
}

func (s *storage) WriteMetric(name string, value any) error {
	s.data[name] = value

	return nil
}

func (s *storage) ReadMetric(name string) (any, bool) {
	value, ok := s.data[name]

	return value, ok
}

func New() MemStorage {
	return &storage{
		data: map[string]any{},
	}
}
