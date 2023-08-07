// Package client provides an http Client to send metrics for server
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"memtracker/internal/config/agent"
	"memtracker/internal/crypt"
	"memtracker/internal/metrics"
	"net/http"

	"github.com/net-byte/go-gateway"
)

const contentTypeKey = "Content-Type"
const contentTypeValue = "application/json"

const acceptEncodingKey = "Accept-Encoding"
const acceptEncodingValue = "gzip"

const xRealIPKey = "X-Real-IP"

const updatesURL = "/updates/"

type Client struct {
	Host        string
	ContentType string
	http.Client
	workers int
	channel chan metrics.Metricable
}

func NewClient() Client {
	log.Println("Running http client")
	wrksCnt := agent.ClientCfg.Limit

	c := Client{
		Host:    agent.ClientCfg.Address,
		Client:  http.Client{},
		workers: wrksCnt,
		channel: make(chan metrics.Metricable, wrksCnt),
	}
	go c.listen()
	return c
}

// listen starts Client instance to work
func (c Client) listen() {
	for i := 0; i < c.workers; i++ {
		go c.work()
	}
}

// Send sends metrics to the server
//
// Pre-cond: listen() method should be called before Send
//
// Post-cond: metrics was send to the server
func (c Client) Send(metrics []metrics.Metricable) {
	for _, metric := range metrics {
		c.channel <- metric
	}
}

func (c Client) send(url string, js []byte) {
	buffer := bytes.NewBuffer(js)
	r, err := c.buildRequest(url, buffer)
	if err != nil {
		log.Println(err)
		return
	}
	c.processRequest(r)
}

func (c Client) buildGauges(metric metrics.Metricable, gauges map[string]interface{}) []metrics.Metrics {
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

func (c Client) buildCounters(metric metrics.Metricable, counter map[string]interface{}) []metrics.Metrics {
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

func (c Client) buildRequest(url string, body io.Reader) (*http.Request, error) {
	request, _ := http.NewRequest(http.MethodPost, url, body)
	request.Header.Add(acceptEncodingKey, acceptEncodingValue)
	request.Header.Add(contentTypeKey, contentTypeValue)

	ip, _ := gateway.DiscoverGatewayIPv4()
	request.Header.Add(xRealIPKey, ip.String())
	return request, nil
}

func (c Client) processRequest(r *http.Request) {
	resp, err := c.Client.Do(r)
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

func (c Client) work() {

	url := agent.ClientCfg.Protocol + c.Host + updatesURL

	for {
		m, ok := <-c.channel
		if !ok {
			return
		}
		c.buildSend(url, m)
	}
}

func (c Client) buildSend(url string, m metrics.Metricable) {
	switch {
	case m.String() == "counter":
		js, err := c.buildMarshalCounters(m)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.send(url, js)
	case m.String() == "gauge":
		js, err := c.buildMarshalGauges(m)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.send(url, js)
	}
}

func (c Client) buildMarshalCounters(m metrics.Metricable) ([]byte, error) {
	toMarshal := c.buildCounters(m, m.AsMap())
	return json.Marshal(toMarshal)
}

func (c Client) buildMarshalGauges(m metrics.Metricable) ([]byte, error) {
	toMarshal := c.buildGauges(m, m.AsMap())
	return json.Marshal(toMarshal)
}
