package metricscollector

import (
	"fmt"
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
	if currentValueRaw, ok := c.storage.ReadMetric(`counter`, name); ok {
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
	metricValue, ok := c.storage.ReadMetric(mtype, name)
	if !ok {
		return nil, fmt.Errorf("unkonwn metric %s", name)
	}

	return metricValue, nil
}

//type count map[string]int64
//type gauge map[string]float64
//
//func (c *metricsCollector) writeGaugeBatch(batch []service.Metrics) ([]service.Metrics, error) {
//	for _, mdata := range batch {
//
//	}
//	return nil, nil
//}

func (c *metricsCollector) SetBatch(batch []Metrics) error {
	//var counters []service.Metrics
	//var gauges []service.Metrics
	//var errs []error

	//counters := ``
	//gauges := ``
	fmt.Printf("\n \n BATCH %#v \n", batch)

	for _, mdata := range batch {
		switch mdata.MType {
		case `counter`:
			updatedValue, err := c.SetCount(mdata.ID, *mdata.Delta)
			if err != nil {
				fmt.Printf("%e", fmt.Errorf("cannot write gauge metric %w", err))
				//errs = append(errs, fmt.Errorf("cannot write gauge metric %w", err))
			}
			mdata.Delta = &updatedValue
			//counters += `(` + mdata.ID + `)`
			//counters = append(counters, mdata)
			//counters[mdata.ID] = *mdata.Delta
		case `gauge`:
			updatedValue, err := c.SetGauge(mdata.ID, *mdata.Value)
			if err != nil {
				fmt.Printf("%e", fmt.Errorf("cannot write gauge metric %w", err))
				//errs = append(errs, fmt.Errorf("cannot write gauge metric %w", err))
			}
			mdata.Value = &updatedValue
			//gauges = append(gauges, mdata)
			//gauges[mdata.ID] = *mdata.Value
		}
	}

	return nil
}
