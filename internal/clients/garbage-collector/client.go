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
	serviceUrl string
}

func New(serviceUrl string) Client {
	return &client{
		serviceUrl: serviceUrl,
	}
}

func (c *client) Gauge(metricName string, metricValue float64) error {
	req, err := http.NewRequest(
		http.MethodPost,
		c.serviceUrl+`/update/gauge/`+metricName+`/`+fmt.Sprintf("%f", metricValue),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-type", "text/plain")

	// пока тело ответа нам не нужно
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Counter(metricName string, metricValue int64) error {
	req, err := http.NewRequest(
		http.MethodPost,
		c.serviceUrl+`/update/counter/`+metricName+`/`+fmt.Sprintf("%d", metricValue),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-type", "text/plain")

	// пока тело ответа нам не нужно
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
