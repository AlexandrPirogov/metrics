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
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestServer() *httptest.Server {
	handler := &DefaultHandler{DB: &db.DB{Storage: db.MemStoageDB()}}
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Post("/update/", handler.UpdateHandler)
	r.Post("/value/", handler.RetrieveMetric)
	return httptest.NewServer(r)

}

func executeUpdateRequest(ts *httptest.Server, body []byte) (*http.Response, error) {
	// Read about http unit testing to eliminate double "application/json"
	return http.Post(ts.URL+"/update/", "application/json", bytes.NewBuffer(body))
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
	values := []float64{-1.123456789, 0.9087865431, 1.112233445, 1.123456789}
	data := []Payload{}
	for i := range values {
		data = append(data, Payload{
			StatusCode: http.StatusCreated,
			Metric: metrics.Metrics{
				ID:    "Some val",
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

			postResp, err := executeUpdateRequest(server, js)
			assert.Nil(t, err)
			postBody, _ := io.ReadAll(postResp.Body)
			var postJsBot metrics.Metrics

			err = json.Unmarshal(postBody, &postJsBot)
			defer postResp.Body.Close()
			assert.Nil(t, err)

			getResp := executeGetValueRequest(server, js)
			defer getResp.Body.Close()

			body, err := io.ReadAll(getResp.Body)
			if err != nil {
				t.Errorf("error while reading resp %v", err)
			}
			var actual metrics.Metrics
			log.Println("Got response starts marshal")
			err = json.Unmarshal(body, &actual)
			if err != nil {
				t.Errorf("error while unmarshal %v", err)
			}
			assert.NotNil(t, actual.Value)
			assert.Nil(t, actual.Delta)
			assert.EqualValues(t, *expected.Metric.Value, *actual.Value)
		})
	}
}

func TestCorrectCounterMetric(t *testing.T) {
	deltas := []int64{1, 2, 3, 4, 5, 6, 7, 8}
	var expectedCounter int64 = 0
	data := []Payload{}
	for i := range deltas {
		data = append(data, Payload{
			StatusCode: http.StatusOK,
			Metric: metrics.Metrics{
				ID:    "PollCount",
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

			resp, err := executeUpdateRequest(server, js)
			require.Nil(t, err)
			defer resp.Body.Close()

			var respJs metrics.Metrics
			buffer, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("got error while reading body %v", err)
			}
			err = json.Unmarshal(buffer, &respJs)
			assert.Nil(t, err)

		})
	}
	p := Payload{
		StatusCode: http.StatusCreated,
		Metric: metrics.Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: nil,
			Value: nil,
		},
	}

	bytes, _ := json.Marshal(p.Metric)
	resp := executeGetValueRequest(server, bytes)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("%v", err)
	}

	var actual metrics.Metrics
	err = json.Unmarshal(body, &actual)
	if err != nil {
		t.Errorf("%v", err)
	}
	assert.Equal(t, expectedCounter, *actual.Delta)
}

func TestCorrectGaugeUpdateHandler(t *testing.T) {
	values := []float64{-1.32198, 0, 1.1, 1.12345}
	data := []Payload{}
	for i := range values {
		data = append(data, Payload{
			StatusCode: http.StatusOK,
			Metric: metrics.Metrics{
				ID:    "some",
				MType: "gauge",
				Delta: nil,
				Value: &values[i],
			},
		})
	}
	//qwes
	server := createTestServer()
	defer server.Close()
	for _, actual := range data {
		t.Run("Correct gauge", func(t *testing.T) {
			js, err := json.Marshal(actual.Metric)
			if err != nil {
				t.Errorf("got error while marshal json %v", err)
			}

			resp, err := executeUpdateRequest(server, js)
			require.Nil(t, err)
			var respJs metrics.Metrics
			buffer, err := io.ReadAll(resp.Body)

			err = json.Unmarshal(buffer, &respJs)

			assert.Nil(t, err)
			assert.EqualValues(t, actual.StatusCode, resp.StatusCode)
			assert.Greater(t, resp.ContentLength, int64(0))
			assert.EqualValues(t, *actual.Metric.Value, *respJs.Value)
			assert.Nil(t, actual.Metric.Delta)
			defer resp.Body.Close()
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

			resp, err := executeUpdateRequest(server, js)
			require.Nil(t, err)
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
			StatusCode: http.StatusOK,
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

			resp, err := executeUpdateRequest(server, js)
			require.Nil(t, err)
			defer resp.Body.Close()

			var respJs metrics.Metrics
			buffer, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("got error while reading body %v", err)
			}
			err = json.Unmarshal(buffer, &respJs)
			assert.Nil(t, err)

			assert.EqualValues(t, actual.StatusCode, resp.StatusCode)
			assert.Greater(t, resp.ContentLength, int64(0))
			require.EqualValues(t, updatedMetric, *respJs.Delta)
			assert.Equal(t, true, true)

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

			resp, err := executeUpdateRequest(server, js)
			require.Nil(t, err)
			defer resp.Body.Close()

			assert.EqualValues(t, actual.StatusCode, resp.StatusCode)
		})
	}
}
