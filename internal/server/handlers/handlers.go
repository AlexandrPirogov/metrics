package handlers

import (
	"log"
	"memtracker/internal/memtrack/metrics"
	"memtracker/internal/server/db"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Handler func
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		params := strings.Split(r.URL.Path, "/")
		if len(params) >= 5 {
			if isUpdatePathCorrect(params) == 0 {
				mtype := params[2]
				mname := params[3]
				val := params[4]
				db.Write(mtype, mname, val)
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(""))
			} else {
				log.Printf("Code: %v\n", isUpdatePathCorrect(params))
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func isUpdatePathCorrect(path []string) int {
	log.Printf("path: %v\n", path)
	metricType := path[2]
	metricName := path[3]
	metricVal := path[4]

	var gauges metrics.MemStats = metrics.MemStats{}
	var counters metrics.Polls = metrics.Polls{}

	var mTypes map[string]bool = make(map[string]bool)
	mTypes[gauges.String()] = true
	mTypes[counters.String()] = true

	// If given incorrect path
	if _, ok := mTypes[metricType]; !ok {
		log.Printf("Given metric type %s not exists!", path[1])
		return 1
	}

	if metricsCheck(&gauges, metricName, metricVal) == 0 || metricsCheck(&counters, metricName, metricVal) == 0 {
		return 0
	}
	return -1
}

func metricsCheck(metric metrics.Metricable, metricName string, metricVal string) int {
	var metrics map[string]bool = make(map[string]bool)
	counterVal := reflect.TypeOf(metric).Elem()
	for i := 0; i < counterVal.NumField(); i++ {
		metrics[counterVal.Field(i).Name] = true
	}

	if _, ok := metrics[metricName]; !ok {
		log.Printf("Metric %s isn't exists!\n", metricName)
		return 1
	}

	if _, err := strconv.ParseFloat(metricVal, 64); err != nil {
		log.Printf("Error: %v", err)
		return 2
	}
	return 0
}
