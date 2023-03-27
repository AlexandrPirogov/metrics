package api

import (
	"encoding/json"
	"fmt"
	"log"
	"memtracker/internal/memtrack/metrics"
	"net/http"
)

type MetricsStorer interface {
	// Reads all metrics and returns their string representation
	Read() []byte
	// Read metrics with given type and name.
	//
	// Pre-cond: Given correct mtype and mname
	//
	// Post-cond: Returns suitable metrics according to given mtype and mname
	ReadByParams(mtype string, mname string) ([]byte, error)
	// Writes metric in store
	//
	// Pre-cond: given correct type name and value of metric
	//
	// Post-cond: stores metric in storage. If success error equals nil
	Write(mtype string, mname string, val string) error
}

type MetricsHandler interface {
	RetrieveMetrics(w http.ResponseWriter, r *http.Request)
	RetrieveMetric(w http.ResponseWriter, r *http.Request)
	UpdateHandler(w http.ResponseWriter, r *http.Request)
}

type DefaultHandler struct {
	DB MetricsStorer
}

// RetrieveMetric returns one metric by given type and name
func (d *DefaultHandler) RetrieveMetric(w http.ResponseWriter, r *http.Request) {
	var metric metrics.Metrics
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil || metric.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
	} else {

		if metric.MType == "gauge" {
			if metric.Value == nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {

				res, err := d.DB.ReadByParams(metric.MType, metric.ID)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					log.Println(err)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(res))
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("{\"a\":\"b\"}"))
				}

			}
		} else if metric.MType == "counter" {
			if metric.Delta == nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				d.DB.Write(metric.MType, metric.ID, fmt.Sprintf("%d", *metric.Delta))
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("{\"a\":\"b\"}"))
			}

		} else {
			w.WriteHeader(http.StatusNotImplemented)
		}

	}

}

// UpdateHandler saves incoming metrics
//
// Pre-cond: given correct type, name and val of metrics
//
// Post-cond: correct metrics saved on server
func (d *DefaultHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	var metric metrics.Metrics
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil || metric.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
	} else {

		if metric.MType == "gauge" {
			if metric.Value == nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				log.Printf("got type:%s  name: %s ID, val: %s", metric.MType, metric.ID, fmt.Sprintf("%f", *metric.Value))
				d.DB.Write(metric.MType, metric.ID, fmt.Sprintf("%f", *metric.Value))
				js, _ := d.DB.ReadByParams(metric.MType, metric.ID)
				log.Printf("got :%s", js)
				w.WriteHeader(http.StatusCreated)
				w.Write(js)
			}
		} else if metric.MType == "counter" {
			if metric.Delta == nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				log.Printf("got type:%s  name: %s ID, val: %s", metric.MType, metric.ID, fmt.Sprintf("%d", *metric.Delta))
				d.DB.Write(metric.MType, metric.ID, fmt.Sprintf("%d", *metric.Delta))
				js, _ := d.DB.ReadByParams(metric.MType, metric.ID)
				log.Printf("got :%s", js)
				w.WriteHeader(http.StatusCreated)
				w.Write(js)
			}

		} else {
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}
