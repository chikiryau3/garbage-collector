package main

import (
	"github.com/chikiryau3/garbage-collector/internal/configs"
	"github.com/chikiryau3/garbage-collector/internal/memStorage"
	"github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	service2 "github.com/chikiryau3/garbage-collector/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := *logger.Sugar()

	config := configs.LoadServiceConfig()

	storage := memstorage.New()
	collector := metricscollector.New(storage)
	service := service2.New(collector, log)

	router := chi.NewRouter()
	//router.Use(middleware.Logger)
	//router.Use(service.WithLogging)

	router.Route(`/update`, func(r chi.Router) {
		r.Post(`/`, service.UpdateHandler)
		r.Route(`/gauge`, func(r chi.Router) {
			r.Route(`/{metricName}/{metricValue}`, func(r chi.Router) {
				r.Post(`/`, service.GaugeHandler)
			})
		})

		r.Route(`/counter`, func(r chi.Router) {
			r.Route(`/{metricName}/{metricValue}`, func(r chi.Router) {
				r.Post(`/`, service.CounterHandler)
			})
		})

		r.Post(`/*`, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		})
	})

	router.Route(`/value`, func(r chi.Router) {
		r.Post(`/`, service.ValueHandler)
		r.Route(`/{metricType}/{metricName}`, func(r chi.Router) {
			r.Get(`/`, service.GetMetric)
		})
	})

	router.Get(`/`, service.GetMetricsHTML)

	err = http.ListenAndServe(config.Endpoint, router)

	if err != nil {
		panic(err)
	}
}
