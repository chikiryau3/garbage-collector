package main

import (
	"github.com/chikiryau3/garbage-collector/internal/configs"
	"github.com/chikiryau3/garbage-collector/internal/logger"
	"github.com/chikiryau3/garbage-collector/internal/memStorage"
	"github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	service2 "github.com/chikiryau3/garbage-collector/internal/service"
	"github.com/chikiryau3/garbage-collector/internal/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	log, err := logger.InitLogger()
	if err != nil {
		panic(err)
	}

	config := configs.LoadServiceConfig()
	storage := memstorage.New(config.StorageConfig)

	if config.Restore {
		err = storage.RestoreFromDump()
		if err != nil {
			log.Error("restore from dump error", err)
		}
	}

	if config.FileStoragePath != "" {
		errs := storage.RunStorageDumper()
		go utils.ListenForErrors(errs, "storage dumper error", log.Error)
	}

	collector := metricscollector.New(storage)
	service := service2.New(collector, log)

	router := chi.NewRouter()
	router.Use(service.WithLogging)
	router.Use(service2.GzipMiddleware)

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
