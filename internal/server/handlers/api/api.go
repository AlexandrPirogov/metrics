package api

import (
	"encoding/json"
	"fmt"
	"log"

	"memtracker/internal/config/server"
	"memtracker/internal/kernel"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/memtrack/metrics"
	"memtracker/internal/server/db"
	"memtracker/internal/server/db/journal"
	"memtracker/internal/server/db/storage/sql/postgres"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func NewHandler() *DefaultHandler {
	if server.ServerCfg.DBUrl != "" {
		err := postgres.Ping()
		if err == nil {
			log.Printf("Using postgres as DB")
			return &DefaultHandler{
				DB: db.DB{
					Storage:   postgres.NewPg(),
					Journaler: journal.NewJournal(),
				},
			}
		}
		log.Printf("err %v", err)

	}

	log.Printf("Using local ram as DB")
	return &DefaultHandler{
		DB: db.DB{
			Storage:   db.MemStorageDB(),
			Journaler: journal.NewJournal(),
		},
	}

}

type DefaultHandler struct {
	DB db.DB
}

// RetrieveMetric return all contained metrics
func (d *DefaultHandler) RetrieveMetrics(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	query := tuples.NewTuple()
	query.SetField("name", "*")
	query.SetField("type", "*")
	res, _ := kernel.Read(d.DB.Storage, query)
	log.Printf("res afte read%v", res)
	body := []byte{}

	for _, tuple := range res {
		b, _ := json.Marshal(tuple)
		body = append(body, b...)
	}
	w.Write(body)

}

// RetrieveMetric returns one metric by given type and name
func (d *DefaultHandler) RetrieveMetric(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	if mtype == "" || mname == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")

	query := tuples.NewTuple()
	query.SetField("name", mname)
	query.SetField("type", mtype)
	res, _ := kernel.Read(d.DB.Storage, query)

	if len(res) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for _, tuple := range res {
		if mtype == "gauge" {
			val, _ := tuple.GetField("value")
			valStr := val.(*float64)
			valB := strconv.FormatFloat(*valStr, 'f', -3, 64)
			w.Write([]byte(valB))
		} else {
			val, _ := tuple.GetField("value")
			valStr := val.(*int64)
			valB := fmt.Sprintf("%d", *valStr)
			w.Write([]byte(valB))
		}
	}

}

// UpdateHandler saves incoming metrics
//
// Pre-cond: given correct type, name and val of metrics
//
// Post-cond: correct metrics saved on server
func (d *DefaultHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	val := chi.URLParam(r, "val")
	if mtype == "" || mname == "" || val == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	code := isUpdatePathCorrect(mtype, mname, val)
	if code != http.StatusOK {
		w.WriteHeader(code)
		return
	}
	metricState, err := metrics.CreateState(mname, mtype, val)
	if err != nil {
		switch err.Error() {
		case "nil value":
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
		return
	}

	_, err = kernel.Write(d.DB.Storage, []tuples.Tupler{metricState})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateHandler saves incoming metrics
//
// Pre-cond: given correct type, name and val of metrics
//
// Post-cond: correct metrics saved on server
func (d *DefaultHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("db url %v", server.ServerCfg)
	err := postgres.Ping()
	log.Printf("ping err %v", err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
		return http.StatusNotImplemented
	}

	if _, err := strconv.ParseFloat(mval, 64); err != nil {
		return http.StatusBadRequest
	}
	return http.StatusOK
}
