package agent

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

// хак с фильтрацией ненужных полей (первое что в голову пришло)
// не уверен в том, что это нормально перформит, но не хотелось вручную все метрики разбирать
func filterFields(data interface{}) (map[string]any, error) {
	var stats RuntimeMetrics
	statsJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(statsJSON, &stats)
	if err != nil {
		return nil, err
	}

	var filteredMap map[string]any
	statsJSON, err = json.Marshal(stats)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(statsJSON, &filteredMap)
	if err != nil {
		return nil, err
	}

	return filteredMap, nil
}

func (a *agent) pollMetrics() error {
	var statsRaw runtime.MemStats
	runtime.ReadMemStats(&statsRaw)

	stats, err := filterFields(statsRaw)
	if err != nil {
		return fmt.Errorf("filter fields error %w", err)
	}

	for metricName, metricValueRaw := range stats {
		metricValue, ok := metricValueRaw.(float64)
		if !ok {
			return fmt.Errorf("metric value type error %w", err)
		}
		// сохраняю метрики с тегами чтобы можно было отправить все в цикле не разбирая каждую отдельно
		_, err := a.collector.SetGauge(`gauge:`+metricName, metricValue)
		if err != nil {
			return fmt.Errorf("set gauge error %w", err)
		}
	}

	// сохраняю метрики с тегами чтобы можно было отправить все в цикле не разбирая каждую отдельно
	_, err = a.collector.SetCount("count:PollCount", 1)
	if err != nil {
		return fmt.Errorf("set count error %w", err)
	}

	// сохраняю метрики с тегами чтобы можно было отправить все в цикле не разбирая каждую отдельно
	_, err = a.collector.SetGauge("gauge:RandomValue", rand.Float64())
	if err != nil {
		return fmt.Errorf("set count error %w", err)
	}

	return nil
}

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

func (a *agent) sendReport() error {
	collectedData, err := a.collector.ReadStorage()
	if err != nil {
		return fmt.Errorf("read storage error %w", err)
	}

	for metricName, metricValueRaw := range *collectedData {
		parts := strings.Split(metricName, `:`)
		metricType := parts[0]
		metricName = parts[1]

		if metricType == `gauge` {
			metricValue, ok := metricValueRaw.(float64)
			if !ok {
				return fmt.Errorf("metric %s has wrong type, value %v", metricName, metricValueRaw)
			}
			err := a.collectionServiceClient.SendGauge(metricName, metricValue)
			if err != nil {
				return fmt.Errorf("send gauge error %w", err)
			}
		}

		if metricType == `count` {
			metricValue, ok := metricValueRaw.(int64)
			if !ok {
				return fmt.Errorf("metric %s has wrong type, value %v", metricName, metricValueRaw)
			}

			err := a.collectionServiceClient.SendCounter(metricName, metricValue)
			if err != nil {
				return fmt.Errorf("send counter error %w", err)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (a *agent) RunReporter() <-chan error {
	errs := make(chan error, 1)
	ticker := time.NewTicker(a.config.ReportInterval)

	go func() {
		for range ticker.C {
			err := a.sendReport()
			if err != nil {
				errs <- fmt.Errorf("send report error %w", err)
				return
			}
		}
	}()

	return errs
}
