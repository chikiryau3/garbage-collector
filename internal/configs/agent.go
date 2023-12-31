package configs

import (
	"flag"
	"github.com/chikiryau3/garbage-collector/internal/agent"
	garbagecollector "github.com/chikiryau3/garbage-collector/internal/clients/garbage-collector"
	"os"
	"strconv"
	"time"
)

type AgentCLIArgs struct {
	serverEndpoint *string
	reportInterval *int64
	pollInterval   *int64
	APIKey         *string
}

type AgentConfig struct {
	ServerEndpoint string
	ReportInterval int64
	PollInterval   int64
	APIKey         string
}

type AgentConfigParsed struct {
	AgentConfig           *agent.Config
	CollectorClientConfig *garbagecollector.Config
}

func LoadAgentConfig() AgentConfigParsed {
	args := &AgentCLIArgs{
		serverEndpoint: flag.String("a", "localhost:8080", "service endpoint"),
		reportInterval: flag.Int64("r", 10, "report interval (seconds)"),
		pollInterval:   flag.Int64("p", 2, "poll interval (seconds)"),
		APIKey:         flag.String("k", "", "api key"),
	}

	flag.Parse()

	config := &AgentConfig{}

	if endpoint, ok := os.LookupEnv(`ADDRESS`); ok {
		config.ServerEndpoint = endpoint
	} else {
		config.ServerEndpoint = *args.serverEndpoint
	}

	if key, ok := os.LookupEnv(`KEY`); ok {
		config.APIKey = key
	} else {
		config.APIKey = *args.APIKey
	}

	if pollInterval, ok := os.LookupEnv(`POLL_INTERVAL`); ok {
		pollIntervalParsed, err := strconv.ParseInt(pollInterval, 10, 8)
		if err != nil {
			config.PollInterval = *args.pollInterval
		} else {
			config.PollInterval = pollIntervalParsed
		}
	} else {
		config.PollInterval = *args.pollInterval
	}

	if reportInterval, ok := os.LookupEnv(`REPORT_INTERVAL`); ok {
		reportIntervalParsed, err := strconv.ParseInt(reportInterval, 10, 8)
		if err != nil {
			config.ReportInterval = *args.reportInterval
		} else {
			config.ReportInterval = reportIntervalParsed
		}
	} else {
		config.ReportInterval = *args.reportInterval
	}

	return AgentConfigParsed{
		&agent.Config{
			PollInterval:   time.Second * time.Duration(config.PollInterval),
			ReportInterval: time.Second * time.Duration(config.ReportInterval),
		},
		&garbagecollector.Config{
			ServiceURL: `http://` + config.ServerEndpoint,
			APIKey:     config.APIKey,
		},
	}
}
