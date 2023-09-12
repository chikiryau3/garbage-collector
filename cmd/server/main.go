package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v9"
	"github.com/chikiryau3/garbage-collector/internal/memStorage"
	"github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	service2 "github.com/chikiryau3/garbage-collector/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Args struct {
	endpoint *string `env:"ADDRESS"`
}

var args = &Args{
	endpoint: flag.String("a", "localhost:8080", "service endpoint"),
}

func main() {
	flag.Parse()

	if err := env.Parse(&args); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", args)

	storage := memstorage.New()
	collector := metricscollector.New(storage)
	service := service2.New(collector)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route(`/update`, func(r chi.Router) {
		r.Route(`/gauge`, func(r chi.Router) {
			r.Route(`/{metricName}/{metricValue}`, func(r chi.Router) {
				r.Use(service.WithMetricData)
				r.Post(`/`, service.GaugeHandler)
			})
		})

		r.Route(`/counter`, func(r chi.Router) {
			r.Route(`/{metricName}/{metricValue}`, func(r chi.Router) {
				r.Use(service.WithMetricData)
				r.Post(`/`, service.CounterHandler)
			})
		})

		r.Post(`/*`, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		})
	})

	router.Route(`/value`, func(r chi.Router) {
		r.Route(`/{metricType}/{metricName}`, func(r chi.Router) {
			r.Use(service.WithMetricData)
			r.Get(`/`, service.GetMetric)
		})
	})

	router.Get(`/`, service.GetMetricsHTML)

	err := http.ListenAndServe(*args.endpoint, router)
	if err != nil {
		panic(err)
	}
}
