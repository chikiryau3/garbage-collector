package main

import (
	"github.com/chikiryau3/garbage-collector/internal/agent"
	garbagecollector "github.com/chikiryau3/garbage-collector/internal/clients/garbage-collector"
	"github.com/chikiryau3/garbage-collector/internal/configs"
	"github.com/chikiryau3/garbage-collector/internal/logger"
	memstorage "github.com/chikiryau3/garbage-collector/internal/memStorage"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"github.com/chikiryau3/garbage-collector/internal/utils"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log, err := logger.InitLogger()
	if err != nil {
		panic(err)
	}

	config := configs.LoadAgentConfig()
	log.Infoln("AGENT CONFIG", config)
	log.Infoln("AGENT ENV", os.Environ())

	storage := memstorage.New(memstorage.DefaultConfig)
	collector := metricscollector.New(storage)
	collectionServiceClient := garbagecollector.New(config.CollectorClientConfig)

	metricsAgent := agent.New(
		collector,
		collectionServiceClient,
		config.AgentConfig,
		log,
	)

	pollErrors := metricsAgent.RunPollChron()
	go utils.ListenForErrors(pollErrors, "storage dumper error", log.Error)

	reporterErrors := metricsAgent.RunReporter()
	go utils.ListenForErrors(reporterErrors, "reporter chron error", log.Error)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}
