package configs

import (
	"flag"
	memstorage "github.com/chikiryau3/garbage-collector/internal/memStorage"
	"os"
	"strconv"
	"time"
)

type ServiceConfig struct {
	StoreInterval   int64
	FileStoragePath string
	Restore         bool
	Endpoint        string
	StorageConfig   *memstorage.Config
}

type ServiceCLIArgs struct {
	storeInterval   *int64
	fileStoragePath *string
	restore         *bool
	endpoint        *string
}

func LoadServiceConfig() *ServiceConfig {
	args := &ServiceCLIArgs{
		endpoint:        flag.String("a", "localhost:8080", "service endpoint"),
		fileStoragePath: flag.String("f", "/tmp/metrics-db.json", "filePath"),
		storeInterval:   flag.Int64("i", 300, "store interval (seconds)"),
		restore:         flag.Bool("r", true, "should restore from file (bool)"),
	}

	flag.Parse()

	config := &ServiceConfig{}

	if endpoint, ok := os.LookupEnv(`ADDRESS`); ok {
		config.Endpoint = endpoint
	} else {
		config.Endpoint = *args.endpoint
	}

	if storeInterval, ok := os.LookupEnv(`STORE_INTERVAL`); ok {
		storeIntervalParsed, err := strconv.ParseInt(storeInterval, 10, 8)
		if err != nil {
			config.StoreInterval = *args.storeInterval
		} else {
			config.StoreInterval = storeIntervalParsed
		}
	} else {
		config.StoreInterval = *args.storeInterval
	}

	if fileStoragePath, ok := os.LookupEnv(`FILE_STORAGE_PATH`); ok {
		config.FileStoragePath = fileStoragePath
	} else {
		config.FileStoragePath = *args.fileStoragePath
	}

	if restore, ok := os.LookupEnv(`RESTORE`); ok {
		config.Restore = restore == `true`
	} else {
		config.Restore = *args.restore
	}

	config.StorageConfig = &memstorage.Config{
		FileStoragePath: config.FileStoragePath,
		StoreInterval:   time.Second * time.Duration(config.StoreInterval),
		SyncStore:       config.StoreInterval == 0,
	}

	return config
}
