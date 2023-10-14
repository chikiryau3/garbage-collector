package service

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// ValueHandler берет данные для обновления метрики из тела запроса (json)
func (s *service) ValueHandler(w http.ResponseWriter, r *http.Request) {
	var mdata MetricsRes
	s.log.Infoln("BODY", r.Body)

	if err := ReadJSONBody(r.Body, &mdata); err != nil {
		s.log.Error("ValueHandler body parsing error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.log.Infoln("MDATA", mdata)

	metricValue, err := s.collector.GetMetric(mdata.ID)
	if err != nil {
		s.log.Error("ValueHandler get metric error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if mdata.MType == `gauge` {
		mdata.Value = metricValue
	} else if mdata.MType == `counter` {
		mdata.Delta = metricValue
	}

	err = WriteJSONBody(w, mdata)
	if err != nil {
		s.log.Error("ValueHandler response writing error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// UpdateHandler берет данные для обновления метрики из тела запроса (json)
func (s *service) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var mdata Metrics
	if err := ReadJSONBody(r.Body, &mdata); err != nil {
		s.log.Error("UpdateHandler body parsing error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mdata.MType == `gauge` {
		value, err := s.collector.SetGauge(mdata.ID, *mdata.Value)
		if err != nil {
			s.log.Error("UpdateHandler set gauge error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		mdata.Value = &value
	} else if mdata.MType == `counter` {
		delta, err := s.collector.SetCount(mdata.ID, *mdata.Delta)
		if err != nil {
			s.log.Error("UpdateHandler set counter error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		mdata.Delta = &delta
	}

	err := WriteJSONBody(w, mdata)
	if err != nil {
		s.log.Error("UpdateHandler response writing error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetMetricsHTML отдает текущий сторадж в виде строки
func (s *service) GetMetricsHTML(w http.ResponseWriter, r *http.Request) {
	data, err := s.collector.ReadStorage()
	if err != nil {
		s.log.Error("GetMetricsHTML read storage error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	_, err = w.Write([]byte(fmt.Sprintf("%v", data)))
	if err != nil {
		s.log.Error("GetMetricsHTML response writing error", err)
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

// CounterHandler берет данные для обновления значения из урла
func (s *service) CounterHandler(w http.ResponseWriter, r *http.Request) {
	mdata := extractMetricsData(r)
	metricValue, err := strconv.ParseInt(mdata.value, 10, 64)
	if err != nil {
		s.log.Error("CounterHandler metric value parsing error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.collector.SetCount(mdata.name, metricValue)
	if err != nil {
		s.log.Error("CounterHandler set counter metric error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")

	w.WriteHeader(http.StatusOK)
}

// GaugeHandler берет данные для обновления значения из урла
func (s *service) GaugeHandler(w http.ResponseWriter, r *http.Request) {
	mdata := extractMetricsData(r)
	metricValue, err := strconv.ParseFloat(mdata.value, 64)
	if err != nil {
		s.log.Error("GaugeHandler metric value parsing error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.collector.SetGauge(mdata.name, metricValue)
	if err != nil {
		s.log.Error("GaugeHandler set gauge metric error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")

	w.WriteHeader(http.StatusOK)
}

// GetMetric берет имя метрики из урла
func (s *service) GetMetric(w http.ResponseWriter, r *http.Request) {
	mdata := extractMetricsData(r)

	val, err := s.collector.GetMetric(mdata.name)
	if err != nil {
		s.log.Error("GetMetric error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("%v", val)))
	if err != nil {
		s.log.Error("GetMetric response writing error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
