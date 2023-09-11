package garbagecollector

import (
	"fmt"
	"net/http"
)

type Client interface {
	Gauge(metricName string, metricValue float64) error
	Counter(metricName string, metricValue int64) error
}

type client struct {
	serviceURL string
}

func New(serviceUrl string) Client {
	return &client{
		serviceURL: serviceUrl,
	}
}

func (c *client) Gauge(metricName string, metricValue float64) error {
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

func (c *client) Counter(metricName string, metricValue int64) error {
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
