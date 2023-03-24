package api

import (
	"bytes"
	"encoding/json"
	"log"
	"memtracker/internal/memtrack/metrics"
	"memtracker/internal/server/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func executeUpdateRequest(body []byte) *http.Response {
	handler := &DefaultHandler{DB: &db.DB{Storage: db.MemStoageDB()}}
	r := chi.NewRouter()
	//r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Post("/update/", handler.UpdateHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()
	// Read about http unit testing to eliminate double "application/json"
	resp, err := http.Post(ts.URL+"/update/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil
	}
	return resp
}

type Payload struct {
	StatusCode int
	Metric     metrics.Metrics
}

func TestCorrectGaugeUpdateHandler(t *testing.T) {
	deltas := []int64{-1, 0, 1}
	//values := []float64{-1.1, 0, 1.1}
	data := []Payload{}
	for _, delta := range deltas {
		data = append(data, Payload{
			StatusCode: http.StatusCreated,
			Metric: metrics.Metrics{
				ID:    "some",
				MType: "gauge",
				Delta: &delta,
				Value: nil,
			},
		})
	}

	for _, actual := range data {
		t.Run("Correct gauge", func(t *testing.T) {
			js, err := json.Marshal(actual.Metric)
			if err != nil {
				t.Errorf("got error while marshal json %v", err)
			}

			resp := executeUpdateRequest(js)
			defer resp.Body.Close()
			log.Printf("%v", resp)
			assert.EqualValues(t, actual.StatusCode, resp.StatusCode)
		})
	}
}

func TestIncorrectGaugeUpdateHandler(t *testing.T) {
	deltas := []int64{-1, 0, 1}
	//values := []float64{-1.1, 0, 1.1}
	data := []Payload{
		{
			StatusCode: http.StatusBadRequest,
			Metric: metrics.Metrics{
				//ID:    "1",
				MType: "gauge1",
				Delta: &deltas[0],
				Value: nil,
			},
		},
		{
			StatusCode: http.StatusBadRequest,
			Metric: metrics.Metrics{
				//ID:    "1",
				MType: "",
				Delta: &deltas[0],
				Value: nil,
			},
		},
		{
			StatusCode: http.StatusBadRequest,
			Metric: metrics.Metrics{
				//ID:    "1",
				MType: "gauge",
				Delta: nil,
				Value: nil,
			},
		},
	}

	for _, actual := range data {
		t.Run("Correct gauge", func(t *testing.T) {
			js, err := json.Marshal(actual.Metric)
			if err != nil {
				t.Errorf("got error while marshal json %v", err)
			}

			resp := executeUpdateRequest(js)
			defer resp.Body.Close()

			assert.EqualValues(t, actual.StatusCode, resp.StatusCode)
		})
	}
}
