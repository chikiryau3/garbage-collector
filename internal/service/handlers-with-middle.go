package service

import (
	"context"
	"fmt"
	"net/http"
)

func (s *service) extractMetricData(ctx context.Context) MetricData {
	mdata := ctx.Value(s.metricDataContextKey)

	return mdata.(MetricData)
}

func (s *service) GaugeHandlerOld(w http.ResponseWriter, r *http.Request) {
	metricDataRaw := s.extractMetricData(r.Context())
	metricName, metricValue, err := s.formatGaugeInput(metricDataRaw.name, metricDataRaw.value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.collector.SetGauge(metricName, metricValue)
	if err != nil {
		//fmt.Printf("COLLECTION ERROR %e \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")

	w.WriteHeader(http.StatusOK)
}

func (s *service) CounterHandlerOld(w http.ResponseWriter, r *http.Request) {
	metricDataRaw := s.extractMetricData(r.Context())
	metricName, metricValue, err := s.formatCounterInput(metricDataRaw.name, metricDataRaw.value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.collector.SetCount(metricName, metricValue)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")

	w.WriteHeader(http.StatusOK)
}

func (s *service) GetMetricOld(w http.ResponseWriter, r *http.Request) {
	metricDataRaw := s.extractMetricData(r.Context())

	val, err := s.collector.GetMetric(metricDataRaw.name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("%v", val)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (s *service) GetMetricsHTMLOld(w http.ResponseWriter, r *http.Request) {
	data, err := s.collector.ReadStorage()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("%v", data)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
}