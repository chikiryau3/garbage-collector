package garbagecollector

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/chikiryau3/garbage-collector/internal/service"
	"github.com/chikiryau3/garbage-collector/internal/utils"
	"net/http"
	"time"
)

// Client -- клиент к сервису, контракты ручек, работа с хедерами, вот это все
type Client interface {
	SendGauge(metricName string, metricValue float64) error
	SendCounter(metricName string, metricValue int64) error
}

type client struct {
	serviceURL string
	retry      backoff.BackOff
}

func New(serviceURL string) Client {
	r := &utils.Retry{
		InitInterval:  time.Second,
		RetryTimeout:  time.Minute,
		MaxRetryTimes: 3,
	}

	return &client{
		retry:      r.NewExponentialBackOff(),
		serviceURL: serviceURL,
	}
}

func (c *client) SendGauge(metricName string, metricValue float64) error {
	var mData service.Metrics
	mData.ID = metricName
	mData.MType = `gauge`
	mData.Value = &metricValue

	body, err := json.Marshal(mData)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	if _, err = g.Write(body); err != nil {
		return err
	}
	if err = g.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.serviceURL+`/update/`,
		&buf,
	)

	if err != nil {
		return fmt.Errorf("request build err %w", err)
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Content-encoding", "gzip")
	req.Header.Add("Accept-encoding", "gzip")

	var res *http.Response
	retriable := func() error {
		res, err = http.DefaultClient.Do(req)
		defer res.Body.Close()
		if err != nil {
			return err
		}

		return nil
	}
	err = backoff.Retry(retriable, c.retry)

	if err != nil {
		fmt.Println(fmt.Errorf("do request err %w", err))
		return nil
	}

	return nil
}

func (c *client) SendCounter(metricName string, metricValue int64) error {
	var mData service.Metrics
	mData.ID = metricName
	mData.MType = `counter`
	mData.Delta = &metricValue

	body, err := json.Marshal(mData)
	if err != nil {
		return fmt.Errorf(" %w", err)
	}

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	if _, err = g.Write(body); err != nil {
		return err
	}
	if err = g.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		c.serviceURL+`/update/`,
		&buf,
	)
	if err != nil {
		return fmt.Errorf("request build err %w", err)
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Content-encoding", "gzip")
	req.Header.Add("Accept-encoding", "gzip")

	var res *http.Response
	retriable := func() error {
		res, err = http.DefaultClient.Do(req)
		defer res.Body.Close()

		if err != nil {
			return err
		}

		err = res.Body.Close()
		if err != nil {
			return fmt.Errorf("body close err %w", err)
		}

		return nil
	}
	err = backoff.Retry(retriable, c.retry)

	if err != nil {
		fmt.Println(fmt.Errorf("do request err %w", err))
		return nil
	}

	return nil
}
