package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (s *service) GaugeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// todo: что-то вернуть
		return
	}

	// todo: разобраться нужно ли чекать контент тайп запроса
	//headers := r.Header
	//if val, ok := headers["Content-Type"]; !ok || val != "text/plain" {
	//	return
	//}

	// ---------- url encoding ----------
	config := s.endpoints[`gauge`]
	url := r.URL
	match, ok := config.pathPattern.Match(url.Path)
	if !ok {
		fmt.Printf("BAD URL %s \n", url.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// ---------- url encoding end ----------

	// ---------- params validation ----------
	metricName := match.Params[`metricName`]
	metricValueRaw := match.Params[`metricValue`]
	metricValueParsed, err := strconv.ParseFloat(metricValueRaw, 64)
	if err != nil {
		fmt.Printf("BAD METRIC VALUE %s \n", metricValueRaw)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// ---------- params validation end ----------

	err = s.collector.Gauge(metricName, metricValueParsed)
	if err != nil {
		fmt.Printf("COLLECTION ERROR %e \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")

	w.WriteHeader(http.StatusOK)
}

func (s *service) CounterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// todo: что-то вернуть
		return
	}

	// ---------- url encoding ----------
	config := s.endpoints[`counter`]
	url := r.URL
	match, ok := config.pathPattern.Match(url.Path)
	if !ok {
		fmt.Printf("BAD URL %s \n", url.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// ---------- url encoding end ----------

	// ---------- params validation ----------
	metricName := match.Params[`metricName`]
	metricValueRaw := match.Params[`metricValue`]
	metricValueParsed, err := strconv.ParseInt(metricValueRaw, 10, 64)
	if err != nil {
		fmt.Printf("BAD METRIC VALUE %s \n", metricValueRaw)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// ---------- params validation end ----------

	err = s.collector.Count(metricName, metricValueParsed)
	if err != nil {
		fmt.Printf("COLLECTION ERROR %e \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")

	w.WriteHeader(http.StatusOK)
}
