package garbagecollector

import (
	"fmt"
	"net/http"
)

// Client -- клиент к сервису, контракты ручек, работа с хедерами, вот это все
type Client interface {
	SendGauge(metricName string, metricValue float64) error
	SendCounter(metricName string, metricValue int64) error
}

type client struct {
	serviceURL string
}

func New(serviceURL string) Client {
	return &client{
		serviceURL: serviceURL,
	}
}

func (c *client) SendGauge(metricName string, metricValue float64) error {
	req, err := http.NewRequest(
		http.MethodPost,
		c.serviceURL+`/update/gauge/`+metricName+`/`+fmt.Sprintf("%f", metricValue),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-type", "text/plain")

	// пока тело ответа нам не нужно
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) SendCounter(metricName string, metricValue int64) error {
	req, err := http.NewRequest(
		http.MethodPost,
		c.serviceURL+`/update/counter/`+metricName+`/`+fmt.Sprintf("%d", metricValue),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-type", "text/plain")

	// пока тело ответа нам не нужно
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
