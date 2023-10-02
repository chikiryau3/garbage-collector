package configs

import (
	"flag"
	"os"
)

type ServiceConfig struct {
	Endpoint string
}

type ServiceCLIArgs struct {
	endpoint *string
}

func LoadServiceConfig() *ServiceConfig {
	args := &ServiceCLIArgs{
		endpoint: flag.String("a", "localhost:8080", "service endpoint"),
	}

	flag.Parse()

	if endpoint, ok := os.LookupEnv(`ADDRESS`); ok {
		return &ServiceConfig{
			Endpoint: endpoint,
		}
	}

	return &ServiceConfig{
		Endpoint: *args.endpoint,
	}
}
