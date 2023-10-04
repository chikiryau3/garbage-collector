package memstorage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func (s *storage) dumpStorage() error {
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

	j = append(j, '\n')
	_, err = file.Write(j)
	if err != nil {
		return fmt.Errorf("dump storage error %w", err)
	}

	return nil
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
	var storageData StorageData

	flags := os.O_RDONLY | os.O_CREATE
	file, err := os.OpenFile(s.config.FileStoragePath, flags, 0666)
	if err != nil {
		return fmt.Errorf("restore from dump error %w", err)
	}

	var buf []byte
	_, err = file.Read(buf)
	if err != nil {
		return fmt.Errorf("restore from dump error %w", err)
	}
	if err := json.Unmarshal(buf, &storageData); err != nil {
		return err
	}

	if storageData != nil {
		s.data = storageData
	}

	return nil
}
