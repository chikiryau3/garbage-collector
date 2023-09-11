package main

import (
	"fmt"
	"github.com/chikiryau3/garbage-collector/internal/agent"
	garbagecollector "github.com/chikiryau3/garbage-collector/internal/clients/garbage-collector"
	memstorage "github.com/chikiryau3/garbage-collector/internal/memStorage"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"net/http"
	"time"
)

func main() {
	storage := memstorage.New()
	collector := metricscollector.New(storage)
	collectionServiceClient := garbagecollector.New("http://localhost:8080")

	metricsAgent := agent.New(collector, collectionServiceClient, time.Second*2, time.Second*10)

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
