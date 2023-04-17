package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"memtracker/internal/config/server"
	"memtracker/internal/crypt"
	"memtracker/internal/kernel/tuples"
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
		return
	}

	err = json.Unmarshal(body, &metric)
	if err != nil || metric.ID == "" {
		log.Printf("unmarshal fails %v", metric)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metricState := metric.ToTuple()
	body, status := d.processRetrieve(metricState)
	w.WriteHeader(status)
	if len(body) > 0 {
		w.Write(body)
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

	if metric.MType != "gauge" && metric.MType != "counter" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	states, err := metrics.ConvertToMetrics([]metrics.Metrics{metric})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := d.verifyHash([]metrics.Metrics{metric}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newStates, err := d.processUpdate(states)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err = d.MarshalStates(newStates)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if len(body) > 0 {
		w.Write(body)
	}
}

func (d *DefaultHandler) UpdatesHandlerJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var metricsSlice []metrics.Metrics
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &metricsSlice)

	if err := d.verifyHash(metricsSlice); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	states, err := metrics.ConvertToMetrics(metricsSlice)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newStates, err := d.processUpdate(states)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err = d.MarshalStates(newStates)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if len(body) > 0 {
		w.Write(body)
	}
}

func (d *DefaultHandler) MarshalStates(states []tuples.Tupler) ([]byte, error) {
	body := make([]byte, 0)
	var err error
	if len(states) == 1 {
		body, err = json.Marshal(states[0])
		if err != nil {
			return []byte{}, err
		}
	} else {
		body, err = json.Marshal(states)
		if err != nil {
			return []byte{}, err
		}
	}
	go func() { d.DB.Journaler.Write(body) }()
	return body, nil
}

func (d *DefaultHandler) verifyHash(metrcs []metrics.Metrics) error {
	var err error = nil
	for _, metric := range metrcs {
		switch metric.MType {
		case "counter":
			err = d.verifyCounterHash(metric)
		case "gauge":
			err = d.verifyGaugeHash(metric)
		default:
			err = errors.New("not implemeneted")
		}
	}
	return err
}

func (d *DefaultHandler) verifyCounterHash(m metrics.Metrics) error {
	key := server.ServerCfg.Hash
	if key == "" {
		return nil
	}
	if m.Delta == nil {
		return errors.New("value must exists")
	}

	check := crypt.Hash(fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta), key)
	if m.Hash != check {
		return errors.New("hash are not equals")
	}
	return nil
}

func (d *DefaultHandler) verifyGaugeHash(m metrics.Metrics) error {
	key := server.ServerCfg.Hash
	if key == "" {
		return nil
	}
	if m.Value == nil {
		return errors.New("value must exists")
	}

	check := crypt.Hash(fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value), key)
	if m.Hash != check {
		return errors.New("hash are not equals")
	}
	return nil
}
