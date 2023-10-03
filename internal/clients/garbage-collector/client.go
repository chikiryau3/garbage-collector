package garbagecollector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chikiryau3/garbage-collector/internal/service"
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
	var mData service.Metrics
	mData.ID = metricName
	mData.MType = `gauge`
	mData.Value = &metricValue

	body, err := json.Marshal(mData)
	if err != nil {
		return err
	}

	//var buf bytes.Buffer
	//g := gzip.NewWriter(&buf)
	//if _, err = g.Write(body); err != nil {
	//	return err
	//}
	//if err = g.Close(); err != nil {
	//	return err
	//}

	fmt.Printf("SendGauge body %s\n", body)

	req, err := http.NewRequest(
		http.MethodPost,
		c.serviceURL+`/update/`,
		bytes.NewBuffer(body),
		//&buf,
	)
	//req.Close = true
	if err != nil {
		return fmt.Errorf("request build err %w", err)
	}

	req.Header.Add("Content-type", "application/json")
	//req.Header.Add("Content-encoding", "gzip")

	// пока тело ответа нам не нужно
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request err %w", err)
	}

	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("body close err %w", err)
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

	//fmt.Printf("SendCounter %#v", mData)
	fmt.Printf("SendCounter body %s\n", body)

	req, err := http.NewRequest(
		http.MethodPost,
		c.serviceURL+`/update/`,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("request build err %w", err)
	}
	//req.Close = true

	req.Header.Add("Content-type", "application/json")

	// пока тело ответа нам не нужно
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request err %w", err)
	}

	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("body close err %w", err)
	}

	return nil
}
