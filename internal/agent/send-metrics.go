package agent

import (
	"fmt"
)

func (a *agent) sendMetric(mdata Metric) error {
	fmt.Printf("METRIC %s\n", mdata.name)

	if mdata.mtype == `gauge` {
		metricValue, ok := mdata.value.(float64)
		if !ok {
			return fmt.Errorf("metric %s has wrong type, value %v", mdata.name, metricValue)
		}
		err := a.collectionServiceClient.SendGauge(mdata.name, metricValue)
		if err != nil {
			return fmt.Errorf("send gauge error %w", err)
		}
	}

	if mdata.mtype == `count` {
		metricValue, ok := mdata.value.(int64)
		if !ok {
			return fmt.Errorf("metric %s has wrong type, value %v", mdata.name, metricValue)
		}

		err := a.collectionServiceClient.SendCounter(mdata.name, metricValue)
		if err != nil {
			return fmt.Errorf("send counter error %w", err)
		}
	}

	return nil
}
