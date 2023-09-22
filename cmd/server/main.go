package main

import (
	"flag"
	"github.com/chikiryau3/garbage-collector/internal/memStorage"
	"github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	service2 "github.com/chikiryau3/garbage-collector/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
)

type Config struct {
	endpoint string
}

type CLIArgs struct {
	endpoint *string
}

func loadConfig() *Config {
	args := &CLIArgs{
		endpoint: flag.String("a", "localhost:8080", "service endpoint"),
	}

	flag.Parse()

	if endpoint, ok := os.LookupEnv(`ADDRESS`); ok {
		return &Config{
			endpoint: endpoint,
		}
	}

	return &Config{
		endpoint: *args.endpoint,
	}
}

func main() {
	config := loadConfig()

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
	println(config.endpoint)
	err := http.ListenAndServe(config.endpoint, router)
	if err != nil {
		panic(err)
	}
}
