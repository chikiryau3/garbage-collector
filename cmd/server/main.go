package main

import (
	"github.com/chikiryau3/garbage-collector/internal/configs"
	"github.com/chikiryau3/garbage-collector/internal/memStorage"
	"github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	service2 "github.com/chikiryau3/garbage-collector/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	config := configs.LoadServiceConfig()

	storage := memstorage.New()
	collector := metricscollector.New(storage)
	service := service2.New(collector)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route(`/update`, func(r chi.Router) {
		r.Route(`/gauge`, func(r chi.Router) {
			r.Route(`/{metricName}/{metricValue}`, func(r chi.Router) {
				// мидлварь не будет работать, если ее зарегать до темплейта {metricName}/{metricValue}
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

	err := http.ListenAndServe(config.Endpoint, router)

	if err != nil {
		panic(err)
	}
}
