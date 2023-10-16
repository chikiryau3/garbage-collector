package main

import (
	"context"
	"github.com/chikiryau3/garbage-collector/internal/configs"
	"github.com/chikiryau3/garbage-collector/internal/logger"
	"github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	service2 "github.com/chikiryau3/garbage-collector/internal/service"
	"github.com/chikiryau3/garbage-collector/internal/utils"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func main() {
	ctx := context.Background()
	log, err := logger.InitLogger()
	if err != nil {
		panic(err)
	}

	config := configs.LoadServiceConfig()

	db, err := utils.InitPgConnection(config.DatabaseDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var storage metricscollector.Storage
	if config.DatabaseDSN == "" {
		storage = utils.InitMemStorage(config, log)
	} else {
		storage, err = utils.InitPgStorage(ctx, db)
		if err != nil {
			panic(err)
		}
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

	router.Get(`/ping`, func(w http.ResponseWriter, r *http.Request) {
		err := db.PingContext(ctx)
		if err != nil {
			log.Error("ping db error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	router.Get(`/`, service.GetMetricsHTML)
	router.Post(`/updates/`, service.BatchUpdateHandler)

	err = http.ListenAndServe(config.Endpoint, router)

	if err != nil {
		panic(err)
	}
}
