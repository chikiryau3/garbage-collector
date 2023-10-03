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

	var mdata Metrics
	if err = json.Unmarshal(buf.Bytes(), &mdata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mdata.MType == `gauge` {
		value, err := s.collector.GetMetric(mdata.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		newValue := value.(float64)
		mdata.Value = &newValue
	} else if mdata.MType == `counter` {
		value, err := s.collector.GetMetric(mdata.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		delta := value.(int64)
		mdata.Delta = &delta
	}

	resp, err := json.Marshal(mdata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mdata.MType == `gauge` {
		metricName, metricValue, err := s.formatGaugeInput(mdata.ID, *mdata.Value)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		value, err := s.collector.SetGauge(metricName, metricValue)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		mdata.Value = value
	} else if mdata.MType == `counter` {
		metricName, metricValue, err := s.formatCounterInput(mdata.ID, *mdata.Delta)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		delta, err := s.collector.SetCount(metricName, metricValue)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		mdata.Delta = delta
	}

	resp, err := json.Marshal(mdata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.collector.SetCount(metricName, metricValue)
	if err != nil {
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.collector.SetGauge(metricName, metricValue)
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
