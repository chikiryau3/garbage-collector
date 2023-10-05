package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *service) ValueHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var mdata MetricsRes
	if err = json.Unmarshal(buf.Bytes(), &mdata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricValue, err := s.collector.GetMetric(mdata.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if mdata.MType == `gauge` {
		mdata.Value = metricValue
	} else if mdata.MType == `counter` {
		mdata.Delta = metricValue
	}

	resp, err := json.Marshal(mdata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *service) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var mdata Metrics
	if err = json.Unmarshal(buf.Bytes(), &mdata); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if mdata.MType == `gauge` {
		value, err := s.collector.SetGauge(mdata.ID, *mdata.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		mdata.Value = &value
	} else if mdata.MType == `counter` {
		delta, err := s.collector.SetCount(mdata.ID, *mdata.Delta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		mdata.Delta = &delta
	}

	resp, err := json.Marshal(mdata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *service) GetMetricsHTML(w http.ResponseWriter, r *http.Request) {
	data, err := s.collector.ReadStorage()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	_, err = w.Write([]byte(fmt.Sprintf("%v", data)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func extractMetricsData(r *http.Request) MetricData {
	return MetricData{
		mtype: chi.URLParam(r, `metricType`),
		name:  chi.URLParam(r, `metricName`),
		value: chi.URLParam(r, `metricValue`),
	}
}

func (s *service) CounterHandler(w http.ResponseWriter, r *http.Request) {
	mdata := extractMetricsData(r)
	metricName, metricValue, err := s.formatCounterInput(mdata.name, mdata.value)

	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.collector.SetCount(metricName, metricValue)
	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")

	w.WriteHeader(http.StatusOK)
}

func (s *service) GaugeHandler(w http.ResponseWriter, r *http.Request) {
	mdata := extractMetricsData(r)
	metricName, metricValue, err := s.formatGaugeInput(mdata.name, mdata.value)
	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.collector.SetGauge(metricName, metricValue)
	if err != nil {
		s.log.Error(err)
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
