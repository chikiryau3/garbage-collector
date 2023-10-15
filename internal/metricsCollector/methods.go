package metricscollector

import (
	"fmt"
)

func (c *metricsCollector) SetGauge(name string, value float64) (float64, error) {
	err := c.storage.WriteMetric(name, value)
	if err != nil {
		return 0, fmt.Errorf("cannot write gauge metric %w", err)
	}

	return value, nil
}

func (c *metricsCollector) writeCount(name string, value int64) (int64, error) {
	err := c.storage.WriteMetric(name, value)
	if err != nil {
		return 0, fmt.Errorf("cannot write count metric %w", err)
	}

	return value, nil
}

func (c *metricsCollector) SetCount(name string, value int64) (int64, error) {
	if currentValueRaw, ok := c.storage.ReadMetric(name); ok {
		if currentValue, ok := currentValueRaw.(int64); ok {
			return c.writeCount(name, value+currentValue)
		}
	}

	return c.writeCount(name, value)
}

func (c *metricsCollector) ReadStorage() (*StorageData, error) {
	return c.storage.GetData()
}

func (c *metricsCollector) GetMetric(name string) (any, error) {
	metricValue, ok := c.storage.ReadMetric(name)
	if !ok {
		return nil, fmt.Errorf("unkonwn metric %s", name)
	}

	return metricValue, nil
}
