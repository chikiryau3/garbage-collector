package metricscollector

import (
	"fmt"
	memstorage "github.com/chikiryau3/garbage-collector/internal/memStorage"
)

func (c *metricsCollector) SetGauge(name string, value float64) (float64, error) {
	err := c.storage.WriteMetric(name, value)
	if err != nil {
		return 0, fmt.Errorf("cannot write gauge metric %w", err)
	}

	return value, nil
}

func (c *metricsCollector) SetCount(name string, value int64) (int64, error) {
	currentValueRaw, ok := c.storage.ReadMetric(name)
	// если метрика еще не представлена в storage, то пишем переданное значение
	if !ok {
		err := c.storage.WriteMetric(name, value)
		if err != nil {
			return 0, fmt.Errorf("cannot write count metric %w", err)
		}

		return value, nil
	}

	// если значение метрики в storage неправильное, то подменяем на переданное
	// UPD: не могу избавиться от этого, тк WriteMetric - это "положи в хранилку что-то"
	// что конкретно класть (то есть какой тип) -- это бизнес-логика, она должна быть в этой структуре
	currentValue, ok := currentValueRaw.(int64)
	if !ok {
		err := c.storage.WriteMetric(name, value)
		if err != nil {
			return 0, fmt.Errorf("cannot write count metric %w", err)
		}

		return value, nil
	}

	newValue := value + currentValue
	err := c.storage.WriteMetric(name, newValue)
	if err != nil {
		return 0, fmt.Errorf("cannot write count metric %w", err)
	}

	return newValue, nil
}

func (c *metricsCollector) ReadStorage() (*memstorage.StorageData, error) {
	return c.storage.GetData()
}

func (c *metricsCollector) GetMetric(name string) (any, error) {
	metricValue, ok := c.storage.ReadMetric(name)
	if !ok {
		return nil, fmt.Errorf("unkonwn metric %s", name)
	}

	return metricValue, nil
}
