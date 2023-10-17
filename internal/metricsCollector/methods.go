package metricscollector

import (
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
)

func (c *metricsCollector) SetGauge(name string, value float64) (float64, error) {
	err := c.storage.WriteMetric(`gauge`, name, value)
	if err != nil {
		return 0, fmt.Errorf("cannot write gauge metric %w", err)
	}

	return value, nil
}

func (c *metricsCollector) writeCount(name string, value int64) (int64, error) {
	err := c.storage.WriteMetric(`counter`, name, value)
	if err != nil {
		return 0, fmt.Errorf("cannot write counter metric %w", err)
	}

	return value, nil
}

func (c *metricsCollector) SetCount(name string, value int64) (int64, error) {
	if currentValueRaw, err := c.storage.ReadMetric(`counter`, name); err == nil {
		if currentValue, ok := currentValueRaw.(int64); ok {
			return c.writeCount(name, value+currentValue)
		}
	}

	return c.writeCount(name, value)
}

func (c *metricsCollector) ReadStorage() (*StorageData, error) {
	return c.storage.GetData()
}

func (c *metricsCollector) GetMetric(mtype string, name string) (any, error) {
	var metricValue any
	var err error

	retriable := func() error {
		metricValue, err = c.storage.ReadMetric(mtype, name)
		var sErr *StorageRetryableError
		if errors.As(err, &sErr) {
			return err
		}

		return backoff.Permanent(err)
	}

	err = backoff.Retry(retriable, c.retry)

	if err != nil {
		return nil, err
	}

	return metricValue, nil
}

func (c *metricsCollector) SetBatch(batch []Metrics) (*[]Metrics, error) {
	var errs error
	for _, mdata := range batch {
		switch mdata.MType {
		case `counter`:
			updatedValue, err := c.SetCount(mdata.ID, *mdata.Delta)
			if err != nil {
				_ = errors.Join(errs, fmt.Errorf("cannot write counter metric %w", err))
			}
			mdata.Delta = &updatedValue
		case `gauge`:
			updatedValue, err := c.SetGauge(mdata.ID, *mdata.Value)
			if err != nil {
				fmt.Printf("%e", fmt.Errorf("cannot write gauge metric %w", err))
				_ = errors.Join(errs, fmt.Errorf("cannot write gauge metric %w", err))
			}
			mdata.Value = &updatedValue
		}
	}

	return &batch, errs
}
