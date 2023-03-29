package api

import (
	"encoding/json"
	"io"
	"memtracker/internal/memtrack/metrics"
	"net/http"
)

// RetrieveMetric returns one metric by given type and name
func (d *DefaultHandler) RetrieveMetricJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
func (d *DefaultHandler) UpdateHandlerJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
