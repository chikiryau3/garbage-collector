package memstorage

import (
	"encoding/json"
	"os"
	"time"
)

func (s *storage) dumpStorage() error {
	data, err := s.GetData()
	if err != nil {
		return err
	}
	flags := os.O_WRONLY | os.O_CREATE
	file, err := os.OpenFile(s.config.FileStoragePath, flags, 0666)
	defer file.Close()

	if err != nil {
		return err
	}

	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	j = append(j, '\n')
	_, err = file.Write(j)
	if err != nil {
		return err
	}
}

func (s *storage) RunStorageDumper() error {
	ticker := time.NewTicker(s.config.StoreInterval)

	go func() {
		for range ticker.C {
			err := s.dumpStorage()
			if err != nil {
				//fmt.Errorf("dump storage error %w", err)
				return
			}
		}
	}()

	return nil
}

func (s *storage) RestoreFromDump() error {
	var storageData *StorageData

	data, err := os.ReadFile(s.config.FileStoragePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, storageData); err != nil {
		return err
	}

	s.data = *storageData

	return nil
}
