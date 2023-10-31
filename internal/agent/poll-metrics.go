package agent

import (
	"fmt"
	"math/rand"
	"runtime"
)

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
