package main

import (
	"flag"
	"fmt"
	"github.com/chikiryau3/garbage-collector/internal/agent"
	garbagecollector "github.com/chikiryau3/garbage-collector/internal/clients/garbage-collector"
	memstorage "github.com/chikiryau3/garbage-collector/internal/memStorage"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"net/http"
	"time"
)

type Args struct {
	serverEndpoint *string
	reportInterval *time.Duration
	pollInterval   *time.Duration
}

var args = &Args{
	serverEndpoint: flag.String("a", "localhost:8080", "service endpoint"),
	reportInterval: flag.Duration("r", 10, "report interval (seconds)"),
	pollInterval:   flag.Duration("p", 2, "poll interval (seconds)"),
}

func main() {
	flag.Parse()

	storage := memstorage.New()
	collector := metricscollector.New(storage)
	collectionServiceClient := garbagecollector.New(`http://` + *args.serverEndpoint)

	metricsAgent := agent.New(
		collector,
		collectionServiceClient,
		time.Second*(*args.pollInterval),
		time.Second*(*args.reportInterval),
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
	err = http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
