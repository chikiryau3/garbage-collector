package service

import (
	"fmt"
	"net/http"
)

func (s *service) GaugeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	metricName, metricValue, err := s.FormatGaugeInput(ctx.Value(`metricName`), ctx.Value(`metricValue`))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.collector.Gauge(metricName, metricValue)
	if err != nil {
		//fmt.Printf("COLLECTION ERROR %e \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")

	w.WriteHeader(http.StatusOK)
}

func (s *service) CounterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	metricName, metricValue, err := s.FormatCounterInput(ctx.Value(`metricName`), ctx.Value(`metricValue`))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.collector.Count(metricName, metricValue)
	if err != nil {
		//fmt.Printf("COLLECTION ERROR %e \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")

	w.WriteHeader(http.StatusOK)
}

func (s *service) GetMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := ctx.Value(`metricType`).(string)

	metricName, ok := ctx.Value(`metricName`).(string)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	val, err := s.collector.GetMetric(metricName)
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

func (s *service) GetMetricsHTML(w http.ResponseWriter, r *http.Request) {
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
