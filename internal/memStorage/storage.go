package memstorage

type StorageData map[string]any

// MemStorage -- интерфейс для работы с хранилищем, "читать и писать"
// сейчас это просто мапа в памяти, но когда-нибудь это станет базой
// в таком случае, можно будет реализовать этот же интерфейс, но с логикой для работы с БД
// поэтому это именно интерфейс (ну и чтобы замокать)
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
