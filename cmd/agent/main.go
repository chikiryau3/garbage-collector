package main

import (
	"fmt"
	"github.com/chikiryau3/garbage-collector/internal/agent"
	garbagecollector "github.com/chikiryau3/garbage-collector/internal/clients/garbage-collector"
	"github.com/chikiryau3/garbage-collector/internal/configs"
	memstorage "github.com/chikiryau3/garbage-collector/internal/memStorage"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := configs.LoadAgentConfig()

	storage := memstorage.New()
	collector := metricscollector.New(storage)
	collectionServiceClient := garbagecollector.New(`http://` + config.ServerEndpoint)

	metricsAgent := agent.New(
		collector,
		collectionServiceClient,
		agent.Config{
			PollInterval:   time.Second * time.Duration(config.PollInterval),
			ReportInterval: time.Second * time.Duration(config.ReportInterval),
		},
	)

	err := metricsAgent.RunPollChron()
	if err != nil {
		fmt.Printf("%e", err)
		return
	}

	err = metricsAgent.RunReporter()
	if err != nil {
		fmt.Printf("%e", err)
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}
