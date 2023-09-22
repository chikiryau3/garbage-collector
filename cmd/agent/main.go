package main

import (
	"flag"
	"fmt"
	"github.com/chikiryau3/garbage-collector/internal/agent"
	garbagecollector "github.com/chikiryau3/garbage-collector/internal/clients/garbage-collector"
	memstorage "github.com/chikiryau3/garbage-collector/internal/memStorage"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"net/http"
	"os"
	"strconv"
	"time"
)

type CLIArgs struct {
	serverEndpoint *string
	reportInterval *int64
	pollInterval   *int64
}

type Config struct {
	serverEndpoint string
	reportInterval int64
	pollInterval   int64
}

func loadConfig() *Config {
	args := &CLIArgs{
		serverEndpoint: flag.String("a", "localhost:8080", "service endpoint"),
		reportInterval: flag.Int64("r", 10, "report interval (seconds)"),
		pollInterval:   flag.Int64("p", 2, "poll interval (seconds)"),
	}

	flag.Parse()

	config := &Config{}

	if endpoint, ok := os.LookupEnv(`ADDRESS`); ok {
		config.serverEndpoint = endpoint
	} else {
		config.serverEndpoint = *args.serverEndpoint
	}

	if pollInterval, ok := os.LookupEnv(`POLL_INTERVAL`); ok {
		pollIntervalParsed, err := strconv.ParseInt(pollInterval, 10, 8)
		if err != nil {
			config.pollInterval = *args.pollInterval
		} else {
			config.pollInterval = pollIntervalParsed
		}
	} else {
		config.pollInterval = *args.pollInterval
	}

	if reportInterval, ok := os.LookupEnv(`REPORT_INTERVAL`); ok {
		reportIntervalParsed, err := strconv.ParseInt(reportInterval, 10, 8)
		if err != nil {
			config.reportInterval = *args.reportInterval
		} else {
			config.reportInterval = reportIntervalParsed
		}
	} else {
		config.reportInterval = *args.reportInterval
	}

	return config
}

func main() {
	config := loadConfig()

	storage := memstorage.New()
	collector := metricscollector.New(storage)
	collectionServiceClient := garbagecollector.New(`http://` + config.serverEndpoint)

	metricsAgent := agent.New(
		collector,
		collectionServiceClient,
		time.Second*time.Duration(config.pollInterval),
		time.Second*time.Duration(config.reportInterval),
	)

	err := metricsAgent.RunPollChron()
	if err != nil {
		fmt.Print(fmt.Errorf("%e", err))
		return
	}

	err = metricsAgent.RunReporter()
	if err != nil {
		fmt.Print(fmt.Errorf("%e", err))
		return
	}

	mux := http.NewServeMux()
	err = http.ListenAndServe(`:8081`, mux)
	if err != nil {
		panic(err)
	}
}
