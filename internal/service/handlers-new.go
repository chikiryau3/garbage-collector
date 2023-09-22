package service

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// хендлеры без мидлвари

// очевидный минус такого: реализация хендлера теперь зависит от выбранного фреймворка
// в первом варианте это было в мидлвари

func (s *service) GaugeHandlerNew(w http.ResponseWriter, r *http.Request) {
	mdata := MetricData{
		mtype: chi.URLParam(r, `metricType`),
		name:  chi.URLParam(r, `metricName`),
		value: chi.URLParam(r, `metricValue`),
	}
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

func (s *service) CounterHandlerNew(w http.ResponseWriter, r *http.Request) {
	mdata := MetricData{
		mtype: chi.URLParam(r, `metricType`),
		name:  chi.URLParam(r, `metricName`),
		value: chi.URLParam(r, `metricValue`),
	}
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
