package service

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// хендлеры без мидлвари

// очевидный минус такого: реализация хендлера теперь зависит от выбранного фреймворка
// в первом варианте это было в мидлвари

func extractMetricsData(r *http.Request) MetricData {
	return MetricData{
		mtype: chi.URLParam(r, `metricType`),
		name:  chi.URLParam(r, `metricName`),
		value: chi.URLParam(r, `metricValue`),
	}
}

func (s *service) GaugeHandler(w http.ResponseWriter, r *http.Request) {
	mdata := extractMetricsData(r)
	metricName, metricValue, err := s.formatGaugeInput(mdata.name, mdata.value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.collector.SetGauge(metricName, metricValue)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")

	w.WriteHeader(http.StatusOK)
}

func (s *service) CounterHandler(w http.ResponseWriter, r *http.Request) {
	mdata := extractMetricsData(r)
	metricName, metricValue, err := s.formatCounterInput(mdata.name, mdata.value)

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

func (s *service) GetMetric(w http.ResponseWriter, r *http.Request) {
	mdata := extractMetricsData(r)

	val, err := s.collector.GetMetric(mdata.name)
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
