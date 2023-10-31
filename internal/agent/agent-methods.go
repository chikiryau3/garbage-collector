package agent

import (
	"fmt"
	"strings"
	"time"
)

func (a *agent) RunPollChron() <-chan error {
	errs := make(chan error, 1)
	ticker := time.NewTicker(a.config.PollInterval)

	go func() {
		for range ticker.C {
			err := a.pollMetrics()
			if err != nil {
				errs <- fmt.Errorf("poll metrics error %w", err)
				return
			}
		}
	}()

	return errs
}

type Metric struct {
	mtype string
	name  string
	value any
}

func (a *agent) worker(jobs <-chan Metric, errs chan<- error) {
	for j := range jobs {
		errs <- a.sendMetric(j)
	}
}

func (a *agent) sendReport(errs chan<- error) {
	collectedData, err := a.collector.ReadStorage()
	if err != nil {
		errs <- fmt.Errorf("read storage error %w", err)
		return
	}

	workerTasks := make(chan Metric, len(*collectedData))
	workerErrs := make(chan error, len(*collectedData))

	for w := 1; int64(w) <= a.config.RateLimit; w++ {
		go a.worker(workerTasks, workerErrs)
	}

	for metricName, metricValueRaw := range *collectedData {
		parts := strings.Split(metricName, `:`)
		metricType := parts[0]
		metricName = parts[1]

		workerTasks <- Metric{
			mtype: metricType,
			name:  metricName,
			value: metricValueRaw,
		}
	}

	for c := 1; c <= len(*collectedData); c++ {
		err := <-workerErrs
		if err != nil {
			errs <- fmt.Errorf("reporter error %w", err)
		}
	}
}

func (a *agent) RunReporter() <-chan error {
	errs := make(chan error)
	ticker := time.NewTicker(a.config.ReportInterval)

	go func() {
		for range ticker.C {
			a.sendReport(errs)
		}
	}()

	return errs
}
