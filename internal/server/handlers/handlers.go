package handlers

import (
	"log"
	"memtracker/internal/memtrack/metrics"
	"net/http"
	"strconv"

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
	// Read metrics value with given type and name.
	//
	// Pre-cond: Given correct mtype and mname
	//
	// Post-cond: Returns current metrics value according to given mtype and mname
	ReadValueByParams(mtype string, mname string) ([]byte, error)
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

// RetrieveMetric return all contained metrics
func (d *DefaultHandler) RetrieveMetrics(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	if mtype == "" || mname == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(d.DB.Read()))
	}
}

// RetrieveMetric returns one metric by given type and name
func (d *DefaultHandler) RetrieveMetric(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	if mtype == "" || mname == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		res, err := d.DB.ReadValueByParams(mtype, mname)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		w.Write([]byte(res))
	}
}

// UpdateHandler saves incoming metrics
//
// Pre-cond: given correct type, name and val of metrics
//
// Post-cond: correct metrics saved on server
func (d *DefaultHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	val := chi.URLParam(r, "val")
	log.Printf("type: %s  name:%s  val:%s\n", mtype, mname, val)
	if mtype == "" || mname == "" || val == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		code := isUpdatePathCorrect(mtype, mname, val)
		if code == http.StatusOK {
			if err := d.DB.Write(mtype, mname, val); err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(""))
		} else {
			w.WriteHeader(code)
		}
	}
}

// isUpdatePathCorrect if given metric is correct
//
// Pre-cond: given string type name and val
//
// Post-cond: if metric is correct -- return 0. Otherwise -1
func isUpdatePathCorrect(mtype, mname, mval string) int {
	var gauges = metrics.MemStats{}
	var counters = metrics.Polls{}

	var mTypes = make(map[string]bool)
	mTypes[gauges.String()] = true
	mTypes[counters.String()] = true

	// If given incorrect path
	if _, ok := mTypes[mtype]; !ok {
		log.Printf("Given metric type %s not exists!", mtype)
		return http.StatusNotImplemented
	}

	if _, err := strconv.ParseFloat(mval, 64); err != nil {
		log.Printf("Error: %v", err)
		return http.StatusBadRequest
	}
	return http.StatusOK
}
