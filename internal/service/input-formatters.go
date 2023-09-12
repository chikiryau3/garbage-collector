package service

import (
	"fmt"
	"strconv"
)

// такое можно сделать через дженерики, но пока не стал запариваться

func (s *service) FormatGaugeInput(metricNameRaw any, metricValueRaw any) (string, float64, error) {
	metricName := metricNameRaw.(string)
	metricValueStr, ok := metricValueRaw.(string)
	metricValueParsed, err := strconv.ParseFloat(metricValueStr, 64)

	if err != nil || !ok {
		return ``, 0, fmt.Errorf("BAD INPUT %v %v", metricName, metricValueRaw)
	}

	return metricName, metricValueParsed, nil
}

func (s *service) FormatCounterInput(metricNameRaw any, metricValueRaw any) (string, int64, error) {
	metricName := metricNameRaw.(string)
	metricValueStr, ok := metricValueRaw.(string)
	metricValueParsed, err := strconv.ParseInt(metricValueStr, 10, 64)

	if err != nil || !ok {
		return ``, 0, fmt.Errorf("BAD INPUT %v %v", metricName, metricValueRaw)
	}

	return metricName, metricValueParsed, nil
}
