package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"memtracker/internal/memtrack/metrics"
	"memtracker/internal/server/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestServer() *httptest.Server {
	handler := &DefaultHandler{DB: &db.DB{Storage: db.MemStoageDB()}}
	r := chi.NewRouter()
	//r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Post("/update/", handler.UpdateHandler)
	r.Post("/value/", handler.RetrieveMetric)
	return httptest.NewServer(r)

}

func executeUpdateRequest(ts *httptest.Server, body []byte) *http.Response {
	// Read about http unit testing to eliminate double "application/json"
	resp, err := http.Post(ts.URL+"/update/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil
	}
	return resp
}

func executeGetValueRequest(ts *httptest.Server, body []byte) *http.Response {
	// Read about http unit testing to eliminate double "application/json"
	resp, err := http.Post(ts.URL+"/value/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil
	}
	return resp
}

type Payload struct {
	StatusCode int
	Metric     metrics.Metrics
}

func TestRetrieveGaugeMetric(t *testing.T) {
	values := []float64{-1.1, 0, 1.1, 1.999}
	data := []Payload{}
	for i := range values {
		data = append(data, Payload{
			StatusCode: http.StatusCreated,
			Metric: metrics.Metrics{
				ID:    "some",
				MType: "gauge",
				Delta: nil,
				Value: &values[i],
			},
		})
	}
	server := createTestServer()
	defer server.Close()
	for _, expected := range data {
		t.Run("Correct gauge", func(t *testing.T) {
			js, err := json.Marshal(expected.Metric)
			if err != nil {
				t.Errorf("got error while marshal json %v", err)
			}

			postResp := executeUpdateRequest(server, js)
			defer postResp.Body.Close()

			getResp := executeGetValueRequest(server, js)
			defer getResp.Body.Close()

			body, err := io.ReadAll(getResp.Body)
			if err != nil {
				log.Fatalf("error while reading resp %v", err)
			}
			var actual metrics.Metrics
			err = json.Unmarshal(body, &actual)
			if err != nil {
				log.Fatalf("error while unmarshal %v", err)
			}
			assert.NotNil(t, actual.Value)
		})
	}
}

func TestCounterGaugeMetric(t *testing.T) {
	deltas := []int64{0, 1, 2, 3, 4, 5}
	var expectedCounter int64 = 0
	data := []Payload{}
	for i := range deltas {
		data = append(data, Payload{
			StatusCode: http.StatusCreated,
			Metric: metrics.Metrics{
				ID:    "some",
				MType: "counter",
				Delta: &deltas[i],
				Value: nil,
			},
		})
		expectedCounter += deltas[i]
	}
	server := createTestServer()
	defer server.Close()
	for _, actual := range data {
		t.Run("Correct counter", func(t *testing.T) {
			js, err := json.Marshal(actual.Metric)
			if err != nil {
				t.Errorf("got error while marshal json %v", err)
			}

			resp := executeUpdateRequest(server, js)
			defer resp.Body.Close()
		})
	}
	p := Payload{
		StatusCode: http.StatusCreated,
		Metric: metrics.Metrics{
			ID:    "some",
			MType: "counter",
			Delta: nil,
			Value: nil,
		},
	}

	bytes, _ := json.Marshal(p)
	resp := executeGetValueRequest(server, bytes)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}

	var actual Payload
	err = json.Unmarshal(body, &actual)
	if err != nil {
		log.Fatalf("%v", err)
	}
	assert.Equal(t, expectedCounter, *actual.Metric.Delta)
}

func TestCorrectGaugeUpdateHandler(t *testing.T) {
	values := []float64{-1.1, 0, 1.1, 1.999}
	data := []Payload{}
	for i := range values {
		data = append(data, Payload{
			StatusCode: http.StatusCreated,
			Metric: metrics.Metrics{
				ID:    "some",
				MType: "gauge",
				Delta: nil,
				Value: &values[i],
			},
		})
	}
	server := createTestServer()
	defer server.Close()
	for _, actual := range data {
		t.Run("Correct gauge", func(t *testing.T) {
			js, err := json.Marshal(actual.Metric)
			if err != nil {
				t.Errorf("got error while marshal json %v", err)
			}

			resp := executeUpdateRequest(server, js)
			defer resp.Body.Close()

			var respJs metrics.Metrics
			buffer, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("%v", err)
			}
			err = json.Unmarshal(buffer, &respJs)
			if err != nil {
				log.Fatalf("error while unmarshal response gauge %v", err)
			}
			log.Printf("here %v", respJs)
			assert.EqualValues(t, actual.StatusCode, resp.StatusCode)
			assert.Greater(t, resp.ContentLength, int64(0))
			assert.EqualValues(t, *actual.Metric.Value, *respJs.Value)
		})
	}
}

func TestIncorrectGaugeUpdateHandler(t *testing.T) {
	values := []float64{-1.1, 0, 1.1}
	data := []Payload{
		{
			StatusCode: http.StatusBadRequest,
			Metric: metrics.Metrics{
				ID:    "",
				MType: "gauge",
				Delta: nil,
				Value: &values[0],
			},
		},
		{
			StatusCode: http.StatusNotImplemented,
			Metric: metrics.Metrics{
				ID:    "1",
				MType: "",
				Delta: nil,
				Value: &values[0],
			},
		},
		{
			StatusCode: http.StatusNotImplemented,
			Metric: metrics.Metrics{
				ID:    "1",
				MType: "gauge1",
				Delta: nil,
				Value: &values[0],
			},
		},
		{
			StatusCode: http.StatusBadRequest,
			Metric: metrics.Metrics{
				ID:    "1",
				MType: "gauge",
				Delta: nil,
				Value: nil,
			},
		},
	}

	server := createTestServer()
	defer server.Close()
	for _, actual := range data {
		t.Run("Incorrect gauge", func(t *testing.T) {
			js, err := json.Marshal(actual.Metric)
			if err != nil {
				t.Errorf("got error while marshal json %v", err)
			}

			resp := executeUpdateRequest(server, js)
			defer resp.Body.Close()

			assert.EqualValues(t, actual.StatusCode, resp.StatusCode)
		})
	}
}

func TestCorrectCounterUpdateHandler(t *testing.T) {
	deltas := []int64{0, 1, 2, 3, 4, 5}
	data := []Payload{}
	for i := range deltas {
		data = append(data, Payload{
			StatusCode: http.StatusCreated,
			Metric: metrics.Metrics{
				ID:    "some",
				MType: "counter",
				Delta: &deltas[i],
				Value: nil,
			},
		})
	}
	server := createTestServer()
	defer server.Close()
	var updatedMetric int64 = 0
	for _, actual := range data {
		updatedMetric += *actual.Metric.Delta
		t.Run("Correct counter", func(t *testing.T) {
			js, err := json.Marshal(actual.Metric)
			if err != nil {
				t.Errorf("got error while marshal json %v", err)
			}

			resp := executeUpdateRequest(server, js)
			defer resp.Body.Close()
			log.Printf("%v", resp)

			var respJs metrics.Metrics
			buffer, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("%v", err)
			}

			if err := json.Unmarshal(buffer, &respJs); err != nil {
				log.Fatalf("%v", err)
			} else {
				assert.EqualValues(t, actual.StatusCode, resp.StatusCode)
				assert.Greater(t, resp.ContentLength, int64(0))
				require.EqualValues(t, updatedMetric, *respJs.Delta)
			}
		})
	}
}

func TestIncorrectCounterUpdateHandler(t *testing.T) {
	deltas := []int64{-1, 0, 1}
	data := []Payload{
		{
			StatusCode: http.StatusBadRequest,
			Metric: metrics.Metrics{
				ID:    "",
				MType: "counter",
				Delta: &deltas[0],
				Value: nil,
			},
		},
		{
			StatusCode: http.StatusNotImplemented,
			Metric: metrics.Metrics{
				ID:    "1",
				MType: "",
				Delta: &deltas[0],
				Value: nil,
			},
		},
		{
			StatusCode: http.StatusBadRequest,
			Metric: metrics.Metrics{
				ID:    "1",
				MType: "counter",
				Delta: nil,
				Value: nil,
			},
		},
	}

	server := createTestServer()
	defer server.Close()
	for _, actual := range data {
		t.Run("Correct counter", func(t *testing.T) {
			js, err := json.Marshal(actual.Metric)
			if err != nil {
				t.Errorf("got error while marshal json %v", err)
			}

			resp := executeUpdateRequest(server, js)
			defer resp.Body.Close()

			assert.EqualValues(t, actual.StatusCode, resp.StatusCode)
		})
	}
}
