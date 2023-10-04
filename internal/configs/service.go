package configs

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type ServiceConfig struct {
	StoreInterval   int64
	FileStoragePath string
	Restore         bool
	Endpoint        string
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

	fmt.Printf("ARGS %#v \n", args)

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
		fmt.Printf("STORE INTERVAL FROM ARGS %d", *args.storeInterval)
	}
	fmt.Printf("CONFIG STORE INTERVAL%#v", config.StoreInterval)

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

	return config
}
