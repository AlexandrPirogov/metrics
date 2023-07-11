package client

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"memtracker/internal/memtrack/metrics"
	"net/http"
)

func NewClient(host, contentType string) Client {
	return Client{
		Host:        host,
		ContentType: contentType,
		Client:      http.Client{},
	}
}

type Client struct {
	Host        string
	ContentType string
	http.Client
}

func (c Client) SendCounter(metric metrics.Metricable, mapMetrics map[string]interface{}) {
	log.Printf("sending to host: %s", c.Host)
	for k, v := range mapMetrics {
		val := v.(float64)
		del := int64(val)
		toMarsal := metrics.Metrics{
			ID:    k,
			MType: metric.String(),
			Delta: &del,
		}
		url := "http://" + c.Host + "/update/"

		js, err := json.Marshal(toMarsal)
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		buffer := bytes.NewBuffer(js)
		resp, err := c.Client.Post(url, c.ContentType, buffer)

		if err != nil {
			log.Print(err)
			continue
		}
		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error while readall %v", err)
		}
	}
}

func (c Client) SendGauges(metric metrics.Metricable, mapMetrics map[string]interface{}) {
	log.Printf("sending to host: %s", c.Host)
	for k, v := range mapMetrics {
		val := v.(float64)
		toMarsal := metrics.Metrics{
			ID:    k,
			MType: metric.String(),
			Value: &val,
		}
		url := "http://" + c.Host + "/update/"

		js, err := json.Marshal(toMarsal)
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		buffer := bytes.NewBuffer(js)
		resp, err := c.Client.Post(url, c.ContentType, buffer)

		if err != nil {
			log.Print(err)
			continue
		}
		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error while readall %v", err)
		}
	}
}
