package handlers

import (
	"log"
	"memtracker/internal/memtrack/metrics"
	"memtracker/internal/server/db"
	"net/http"
	"strconv"
	"strings"
)

func NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		params := strings.Split(r.URL.Path, "/")
		if len(params) >= 5 {
			code := isUpdatePathCorrect(params)
			if code == http.StatusOK {
				mtype := params[2]
				mname := params[3]
				val := params[4]
				db.Write(mtype, mname, val)
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(""))
			} else {
				log.Printf("Code: %v\n", isUpdatePathCorrect(params))
				w.WriteHeader(code)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func isUpdatePathCorrect(path []string) int {
	log.Printf("path: %v\n", path)
	metricType := path[2]
	metricName := path[3]
	metricVal := path[4]

	if metricType == "" || metricName == "" || metricVal == "" {
		return http.StatusNotFound
	}
	var gauges = metrics.MemStats{}
	var counters = metrics.Polls{}

	var mTypes = make(map[string]bool)
	mTypes[strings.ToLower(gauges.String())] = true
	mTypes[strings.ToLower(counters.String())] = true

	// If given incorrect path
	if _, ok := mTypes[strings.ToLower(metricType)]; !ok {
		log.Printf("Given metric type %s not exists!", path[1])
		return http.StatusNotImplemented
	}

	if _, err := strconv.ParseFloat(metricVal, 64); err != nil {
		log.Printf("Error: %v", err)
		return http.StatusBadRequest
	}

	return http.StatusOK
}
