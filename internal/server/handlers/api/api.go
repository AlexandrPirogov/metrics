package api

import (
	"encoding/json"
	"io"
	"log"
	"memtracker/internal/memtrack/metrics"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	Write(mtype string, mname string, val string) ([]byte, error)
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		err := json.Unmarshal(body, &metric)
		if err != nil || metric.ID == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			body, status := d.processRetrieve(metric)
			w.WriteHeader(status)
			if len(body) > 0 {
				w.Write(body)
			}
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
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &metric)
	if err != nil || metric.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		body, status := d.processUpdate(metric)
		w.WriteHeader(status)
		if len(body) > 0 {
			w.Write(body)
		}
	}

}

// RetrieveMetric return all contained metrics
func (d *DefaultHandler) RetrieveMetrics(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got by value")
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	if mtype == "" || mname == "" {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(d.DB.Read()))
	}
}
