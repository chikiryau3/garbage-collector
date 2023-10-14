package memstorage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func (s *storage) dumpStorage() error {
	s.Lock()
	defer s.Unlock()
	data, err := s.GetData()
	if err != nil {
		return fmt.Errorf("dump storage error %w", err)
	}
	flags := os.O_WRONLY | os.O_CREATE
	file, err := os.OpenFile(s.config.FileStoragePath, flags, 0666)
	defer func() {
		err = file.Close()
	}()

	if err != nil {
		return fmt.Errorf("dump storage error %w", err)
	}

	j, err := json.Marshal(data)

	if err != nil {
		return fmt.Errorf("dump storage error %w", err)
	}

	_, err = file.Write(j)

	if err != nil {
		return fmt.Errorf("dump storage error %w", err)
	}

	return nil
}

func (s *storage) RunStorageDumper() <-chan error {
	errs := make(chan error, 1)
	ticker := time.NewTicker(s.config.StoreInterval)

	go func() {
		for range ticker.C {
			err := s.dumpStorage()

			if err != nil {
				errs <- fmt.Errorf("dumper error %e", err)
				return
			}
		}
	}()

	return errs
}

func (s *storage) RestoreFromDump() error {
	s.Lock()
	defer s.Unlock()

	var storageData StorageData
	data, err := os.ReadFile(s.config.FileStoragePath)

	if err != nil {
		return fmt.Errorf("restore from dump error %w", err)
	}
	if err := json.Unmarshal(data, &storageData); err != nil {
		return fmt.Errorf("restore from dump error %w", err)
	}

	if storageData != nil {
		s.data = storageData
	}

	return nil
}
