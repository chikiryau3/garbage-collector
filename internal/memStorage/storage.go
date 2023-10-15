package memstorage

import (
	"fmt"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"os"
	"sync"
	"time"
)

// MemStorage -- интерфейс для работы с хранилищем, "читать и писать"
// сейчас это просто мапа в памяти, но когда-нибудь это станет базой
// в таком случае, можно будет реализовать этот же интерфейс, но с логикой для работы с БД
// поэтому это именно интерфейс (ну и чтобы замокать)
type MemStorage metricscollector.Storage

type Config struct {
	FileStoragePath string
	StoreInterval   time.Duration
	SyncStore       bool
}

type storage struct {
	data   metricscollector.StorageData
	config *Config
	sync.Mutex
}

func appendDataToFile(path string, data []byte) error {
	flags := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	file, err := os.OpenFile(path, flags, 0666)
	defer func() {
		err = file.Close()
	}()

	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) WriteMetric(name string, value any) error {
	s.Lock()
	defer s.Unlock()
	s.data[name] = value

	if s.config.SyncStore {
		err := appendDataToFile(s.config.FileStoragePath, []byte(fmt.Sprintf("\"%s\":\"%v\"", name, value)))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *storage) ReadMetric(name string) (any, bool) {
	s.Lock()
	defer s.Unlock()
	value, ok := s.data[name]

	return value, ok
}

func (s *storage) GetData() (*metricscollector.StorageData, error) {
	return &s.data, nil
}

var DefaultConfig = &Config{
	FileStoragePath: "",
	StoreInterval:   0,
	SyncStore:       false,
}

func New(c *Config) MemStorage {
	return &storage{
		data:   map[string]any{},
		config: c,
	}
}
