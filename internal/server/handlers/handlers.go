package handlers

import (
	"log"
	"memtracker/internal/memtrack/metrics"
	"memtracker/internal/server/db"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetDocuments returns saved metrics on server
func GetDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(db.Read()))
}

// UpdateHandler saves incoming metrics
//
// Pre-cond: given correct type, name and val of metrics
//
// Post-cond: correct metrics saved on server
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	val := chi.URLParam(r, "val")
	log.Printf("type: %s  name:%s  val:%s\n", mtype, mname, val)
	if mtype == "" || mname == "" || val == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		code := isUpdatePathCorrect(mtype, mname, val)
		log.Printf("Correct metrics %d\n", code)
		if code == 0 {
			db.Write(mtype, mname, val)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(""))
		} else {
			w.WriteHeader(http.StatusNotFound)
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
		return -1
	}

	if metricsCheck(&gauges, mname, mval) == 0 || metricsCheck(&counters, mname, mval) == 0 {
		return 0
	}
	return -1
}

// metricsCheck checks if given name is correct for given type of metric
//
// Pre-cond: given string type and val
//
// Post-cond: if metric is correct -- return 0. Otherwise -1
func metricsCheck(metric metrics.Metricable, metricName string, metricVal string) int {
	var metrics = make(map[string]bool)
	counterVal := reflect.TypeOf(metric).Elem()
	for i := 0; i < counterVal.NumField(); i++ {
		metrics[counterVal.Field(i).Name] = true
	}

	if _, ok := metrics[metricName]; !ok {
		log.Printf("Metric %s isn't exists!\n", metricName)
		return -1
	}

	if _, err := strconv.ParseFloat(metricVal, 64); err != nil {
		log.Printf("Error: %v", err)
		return -1
	}
	return 0
}
