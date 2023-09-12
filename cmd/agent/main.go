package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v9"
	"github.com/chikiryau3/garbage-collector/internal/agent"
	garbagecollector "github.com/chikiryau3/garbage-collector/internal/clients/garbage-collector"
	memstorage "github.com/chikiryau3/garbage-collector/internal/memStorage"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"net/http"
	"time"
)

type Args struct {
	serverEndpoint *string
	reportInterval *int `env:"REPORT_INTERVAL"`
	pollInterval   *int `env:"POLL_INTERVAL"`
}

var args = &Args{
	serverEndpoint: flag.String("a", "localhost:8080", "service endpoint"),
	reportInterval: flag.Int("r", 10, "report interval (seconds)"),
	pollInterval:   flag.Int("p", 2, "poll interval (seconds)"),
}

func main() {
	flag.Parse()

	if err := env.Parse(&args); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", args)

	storage := memstorage.New()
	collector := metricscollector.New(storage)
	collectionServiceClient := garbagecollector.New(`http://` + *args.serverEndpoint)

	metricsAgent := agent.New(
		collector,
		collectionServiceClient,
		time.Second*time.Duration(*args.pollInterval),
		time.Second*time.Duration(*args.reportInterval),
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
