package handlers

import (
	"log"
	"memtracker/internal/memtrack/metrics"
	"memtracker/internal/server/db"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func RetrieveMetrics(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	if mtype == "" || mname == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(db.Read()))
}

func RetrieveMetric(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	if mtype == "" || mname == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		res, code := db.ReadByParams(mtype, mname)
		w.WriteHeader(code)
		w.Write([]byte(res))
	}
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
		if code == http.StatusOK {
			db.Write(mtype, mname, val)
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
