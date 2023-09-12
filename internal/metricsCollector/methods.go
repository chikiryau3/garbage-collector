package metricscollector

import (
	"fmt"
	memstorage "github.com/chikiryau3/garbage-collector/internal/memStorage"
)

func (c *metricsCollector) Gauge(name string, value float64) error {
	//fmt.Printf("Gauge: %s %f \n", name, value)

	_, ok := c.storage.ReadMetric(name)
	// если метрика еще не представлена в storage
	if !ok {
		err := c.storage.WriteMetric(name, value)
		if err != nil {
			return err
		}

		return nil
	}

	err := c.storage.WriteMetric(name, value)
	if err != nil {
		return err
	}

	return nil
}

func (c *metricsCollector) Count(name string, value int64) error {
	//fmt.Printf("COUNT: %s %d \n", name, value)

	currentValueRaw, ok := c.storage.ReadMetric(name)
	// если метрика еще не представлена в storage, то пишем переданное значение
	if !ok {
		err := c.storage.WriteMetric(name, value)
		if err != nil {
			return err
		}

		return nil
	}

	// если значение метрики в storage неправильное, то подменяем на переданное
	currentValue, ok := currentValueRaw.(int64)
	if !ok {
		err := c.storage.WriteMetric(name, value)
		if err != nil {
			return err
		}

		return nil
	}

	//fmt.Printf("CURRENT VALUE %d", currentValue)

	err := c.storage.WriteMetric(name, value+currentValue)
	if err != nil {
		return err
	}

	return nil
}

func (c *metricsCollector) ReadStorage() (*memstorage.StorageData, error) {
	return c.storage.GetData()
}

func (c *metricsCollector) GetMetric(name string) (any, error) {
	metricValue, ok := c.storage.GetValue(name)
	if !ok {
		return nil, fmt.Errorf("unkonwn metric")
	}

	return metricValue, nil
}
