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
		return err
	}

	//fmt.Printf("runtime %#v", stats)
	//fmt.Print("collect \n")

	for metricName, metricValueRaw := range stats {
		metricValue, ok := metricValueRaw.(float64)
		if !ok {
			return err
		}
		// сохраняю метрики с тегами чтобы можно было отправить все в цикле не разбирая каждую отдельно
		err := a.collector.Gauge(`gauge:`+metricName, metricValue)
		if err != nil {
			return err
		}
	}

	// сохраняю метрики с тегами чтобы можно было отправить все в цикле не разбирая каждую отдельно
	err = a.collector.Count("count:PollCount", 1)
	if err != nil {
		return err
	}

	// сохраняю метрики с тегами чтобы можно было отправить все в цикле не разбирая каждую отдельно
	err = a.collector.Count("count:RandomValue", rand.Int63())
	if err != nil {
		return err
	}

	return nil
}

func (a *agent) RunPollChron() error {
	ticker := time.NewTicker(a.pollInterval)

	go func() {
		for range ticker.C {
			err := a.pollMetrics()
			if err != nil {
				//fmt.Print(fmt.Errorf("poll error %e", err))
				return
			}
		}
	}()

	return nil
}

func (a *agent) sendReport() error {
	collectedData, err := a.collector.ReadStorage()
	if err != nil {
		return err
	}
	//fmt.Print("send \n")

	for metricName, metricValueRaw := range *collectedData {
		parts := strings.Split(metricName, `:`)
		metricType := parts[0]
		metricName = parts[1]

		if metricType == `gauge` {
			metricValue, ok := metricValueRaw.(float64)
			if !ok {
				return fmt.Errorf("metric %s has wrong type, value %T", metricName, metricValueRaw)
			}
			err := a.collectionServiceClient.Gauge(metricName, metricValue)
			if err != nil {
				return err
			}
		}

		if metricType == `count` {
			metricValue, ok := metricValueRaw.(int64)
			if !ok {
				return fmt.Errorf("metric %s has wrong type, value %T", metricName, metricValueRaw)
			}
			err := a.collectionServiceClient.Counter(metricName, metricValue)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *agent) RunReporter() error {
	ticker := time.NewTicker(a.reportInterval)

	go func() {
		for range ticker.C {
			err := a.sendReport()
			if err != nil {
				//fmt.Print(fmt.Errorf("report error %e", err))
				return
			}
		}
	}()

	return nil
}
