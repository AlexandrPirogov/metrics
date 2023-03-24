package api

import (
	"encoding/json"
	"fmt"
	"memtracker/internal/memtrack/metrics"
	"net/http"
)

type MetricsStorer interface {
	// Reads all metrics and returns their string representation
	Read() string
	// Read metrics with given type and name.
	//
	// Pre-cond: Given correct mtype and mname
	//
	// Post-cond: Returns suitable metrics according to given mtype and mname
	ReadByParams(mtype string, mname string) (string, error)
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
				d.DB.Write(metric.MType, metric.ID, fmt.Sprintf("%f", *metric.Value))
				w.WriteHeader(http.StatusCreated)
			}
		} else if metric.MType == "counter" {
			if metric.Delta == nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				d.DB.Write(metric.MType, metric.ID, fmt.Sprintf("%d", *metric.Delta))
				w.WriteHeader(http.StatusCreated)
			}

		} else {
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}
