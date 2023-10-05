package service

import (
	"fmt"
	"strconv"
)

// такое можно сделать через дженерики, но пока не стал запариваться
// UPD: имеется в виду, что это два разных метода только потому, что разные типы метрик
// по сути, это просто приведение типов из строк, в нужные для хранилки
// вынес, чтобы хендлеры не были простыней
// сложности не должно добавить, тк названия функций говорящие (мне кажется)

func (s *service) formatGaugeInput(metricNameRaw any, metricValueRaw any) (string, float64, error) {
	metricName := metricNameRaw.(string)
	metricValueStr, ok := metricValueRaw.(string)

	metricValueParsed, err := strconv.ParseFloat(metricValueStr, 64)

	if err != nil || !ok {
		return ``, 0, fmt.Errorf("BAD INPUT %v %v", metricName, metricValueRaw)
	}

	return metricName, metricValueParsed, nil
}

func (s *service) formatCounterInput(metricNameRaw any, metricValueRaw any) (string, int64, error) {
	metricName := metricNameRaw.(string)
	metricValueStr, ok := metricValueRaw.(string)
	metricValueParsed, err := strconv.ParseInt(metricValueStr, 10, 64)

	if err != nil || !ok {
		return ``, 0, fmt.Errorf("BAD INPUT %v %v", metricName, metricValueRaw)
	}

	return metricName, metricValueParsed, nil
}
