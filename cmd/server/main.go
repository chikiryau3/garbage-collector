package main

import (
	"github.com/chikiryau3/garbage-collector/internal/memStorage"
	"github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"github.com/ucarion/urlpath"
	"net/http"
)

type Service interface {
	GaugeHandler(w http.ResponseWriter, r *http.Request)
	CounterHandler(w http.ResponseWriter, r *http.Request)
}

type service struct {
	collector metricscollector.MetricsCollector
	endpoints map[string]endpoint
}

func New(collector metricscollector.MetricsCollector, endpoints endpoints) Service {
	return &service{
		collector: collector,
		endpoints: endpoints,
	}
}

type endpoint struct {
	path        string
	pathPattern urlpath.Path
}

type endpoints map[string]endpoint

func main() {
	storage := memstorage.New()
	collector := metricscollector.New(storage)
	service := New(collector, map[string]endpoint{
		`gauge`: {
			path:        `/update/gauge/`,
			pathPattern: urlpath.New(`/update/gauge/:metricName/:metricValue`),
		},
		`counter`: {
			path:        `/update/counter/`,
			pathPattern: urlpath.New(`/update/counter/:metricName/:metricValue`),
		},
	})

	mux := http.NewServeMux()
	//mux.HandleFunc(`/`, fn)
	mux.HandleFunc(`/update/gauge/`, service.GaugeHandler)
	mux.HandleFunc(`/update/counter/`, service.CounterHandler)
	mux.HandleFunc(`/update/`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
