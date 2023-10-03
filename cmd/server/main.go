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
	router.Use(service.WithLogging)

	router.Post(`/update`, service.UpdateHandler)

	router.Post(`/value`, service.ValueHandler)

	router.Get(`/`, service.GetMetricsHTML)

	err = http.ListenAndServe(config.Endpoint, router)

	if err != nil {
		panic(err)
	}
}
