package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"memtracker/internal/config/agent"
	"memtracker/internal/crypt"
	"memtracker/internal/memtrack/metrics"
	"net/http"
)

func NewClient(host, contentType string) Client {
	wrksCnt := agent.ClientCfg.Limit
	return Client{
		Host:        host,
		ContentType: contentType,
		Client:      http.Client{},
		workers:     wrksCnt,
		channel:     make(chan metrics.Metricable, wrksCnt),
	}
}

type Client struct {
	Host        string
	ContentType string
	http.Client
	workers int
	channel chan metrics.Metricable
}

func (c Client) Send(metrics []metrics.Metricable) {
	for _, metric := range metrics {
		c.channel <- metric
	}
}

func (c Client) Listen() {
	for i := 0; i < c.workers; i++ {
		go c.work()
	}
}

func (c Client) work() {
	for {
		m, ok := <-c.channel
		if !ok {
			return
		}

		switch {
		case m.String() == "counter":
			c.SendCounter(m, m.AsMap())
		case m.String() == "gauge":
			c.SendGauges(m, m.AsMap())
		}
	}
}

func (c Client) SendCounter(metric metrics.Metricable, mapMetrics map[string]interface{}) {
	toMarshal := c.BuildCounters(metric, mapMetrics)
	url := "http://" + c.Host + "/updates/"
	log.Printf("sending to host: %s", url)

	js, err := json.Marshal(toMarshal)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	buffer := bytes.NewBuffer(js)
	request, _ := http.NewRequest(http.MethodPost, url, buffer)
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Content-Type", "application/json")
	resp, err := c.Client.Do(request)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while readall %v", err)
	}

}

func (c Client) SendGauges(metric metrics.Metricable, mapMetrics map[string]interface{}) {
	toMarshal := c.BuildGauges(metric, mapMetrics)
	url := "http://" + c.Host + "/updates/"
	log.Printf("sending to host: %s", url)

	js, err := json.Marshal(toMarshal)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	buffer := bytes.NewBuffer(js)
	request, _ := http.NewRequest(http.MethodPost, url, buffer)
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Content-Type", "application/json")
	resp, err := c.Client.Do(request)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while readall %v", err)
	}
}

func (c Client) BuildGauges(metric metrics.Metricable, gauges map[string]interface{}) []metrics.Metrics {
	key := agent.ClientCfg.Hash
	res := make([]metrics.Metrics, 0)
	for k, v := range gauges {
		val := v.(float64)
		toMarsal := metrics.Metrics{
			ID:    k,
			MType: metric.String(),
			Value: &val,
			Hash:  crypt.Hash(fmt.Sprintf("%s:gauge:%f", k, val), key),
		}
		res = append(res, toMarsal)
	}
	return res
}

func (c Client) BuildCounters(metric metrics.Metricable, counter map[string]interface{}) []metrics.Metrics {
	key := agent.ClientCfg.Hash
	res := make([]metrics.Metrics, 0)
	for k, v := range counter {
		val := v.(float64)
		del := int64(val)
		toMarsal := metrics.Metrics{
			ID:    k,
			MType: metric.String(),
			Delta: &del,
			Hash:  crypt.Hash(fmt.Sprintf("%s:counter:%d", k, del), key),
		}
		res = append(res, toMarsal)
	}
	return res
}
