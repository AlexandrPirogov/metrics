package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"memtracker/internal/memtrack/metrics"
	"net/http"
)

func ExampleDefaultHandler_RetrieveMetrics() {
	// Setting params for request
	url := "http://" + "your host" + "/value/"
	request, _ := http.NewRequest(http.MethodPost, url, nil)
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Content-Type", "application/json")
	client := http.Client{}

	//Making request
	resp, err := client.Do(request)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	// If success do whatever you want
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while readall %v", err)
	}
}

func ExampleDefaultHandler_RetrieveMetric() {
	// Creating payload pattern to retrieve corresponding metric
	cntr := metrics.Polls{
		PollCount: 1,
	}
	js, err := json.Marshal(cntr)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	// Setting params for request
	url := "http://" + "your host" + "/value/"
	buffer := bytes.NewBuffer(js)
	request, _ := http.NewRequest(http.MethodPost, url, buffer)
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Content-Type", "application/json")
	client := http.Client{}

	// Making request
	resp, err := client.Do(request)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	// If success do whatever you want
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while readall %v", err)
	}
}

func ExampleDefaultHandler_UpdateHandler() {
	// Creating payload pattern to update metric's vallue
	cntr := metrics.Polls{
		PollCount: 1,
	}
	js, err := json.Marshal(cntr)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	// Setting params for request
	url := "http://" + "your host" + "/update/"
	buffer := bytes.NewBuffer(js)
	request, _ := http.NewRequest(http.MethodPost, url, buffer)
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Content-Type", "application/json")
	client := http.Client{}

	// Making request
	resp, err := client.Do(request)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	// If success you'll get an success msg
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while readall %v", err)
	}
}
