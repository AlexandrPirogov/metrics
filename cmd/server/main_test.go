package main_test

import (
	"fmt"
	"io"
	"log"
	"memtracker/internal/memtrack/metrics"
	"memtracker/internal/server/handlers"
	"memtracker/internal/tests"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type response struct {
	code        int
	contentType string
	response    string
}

const host string = "http://localhost:8080"
const path string = "/update"

// #TODO run get method

func TestUpdateHandlerIncorrectPath(t *testing.T) {
	expectFail := response{http.StatusNotFound, "text/plain; charset=utf-8", "404 page not found\n"}
	incorrectPaths := []string{
		path + "/qwe/asd",
		"/qwe",
		"/A/A/A",
		"/A/B/B/B",
		"/C/B/_/E/W/",
	}

	for _, url := range incorrectPaths {
		t.Run(url, func(t *testing.T) {
			//Running server and executing request
			res := runPost(url)

			//Defering to close body
			defer res.Body.Close()

			//Check response code
			tests.AssertEqualValues(t, expectFail.code, res.StatusCode)

			//Trying to read body
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			//Check for Content-type values
			tests.AssertHeader(t, res, "Content-Type", expectFail.contentType)
			//Check for body response
			tests.AssertEqualValues(t, expectFail.response, string(resBody))
		})
	}
}

func TestUpdateHandlerCorrectPath(t *testing.T) {
	expectSucces := response{200, "text/plain", ""}

	correctPaths := CorrectPaths()

	for _, url := range correctPaths {
		t.Run(url, func(t *testing.T) {
			//Running server and executing request
			res := runPost(url)
			//Defering to close body
			defer res.Body.Close()

			//Check response code
			tests.AssertEqualValues(t, expectSucces.code, res.StatusCode)

			//Trying to read body
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			//Check for Content-type values
			tests.AssertHeader(t, res, "Content-Type", expectSucces.contentType)
			//Check for body response
			tests.AssertEqualValues(t, expectSucces.response, string(resBody))
			runGet("value/")
		})
	}
}

func runPost(url string) *http.Response {
	//Running server
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/update/{mtype}/{mname}/{val}", handlers.UpdateHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()
	log.Printf("Url: %v\n", ts.URL+url)
	resp, err := http.Post(ts.URL+url, "text/plain", nil)
	if err != nil {
		return nil
	}
	return resp
}

func runGet(url string) *http.Response {
	//Running server
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/value/{mtype}/{mname}", handlers.UpdateHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()
	log.Printf("Url: %v\n", ts.URL+url)
	resp, err := http.Get(ts.URL + url)
	if err != nil {
		return nil
	}
	return resp
}

func CorrectPaths() []string {
	var gauges = metrics.MemStats{}
	var paths = make([]string, 0)

	gaugeVal := reflect.ValueOf(gauges)
	for i := 0; i < gaugeVal.NumField(); i++ {
		url := fmt.Sprintf("%s/%s/%v/%v", path, gauges, gaugeVal.Field(i).Type().Name(), gaugeVal.Field(i))
		paths = append(paths, url)
	}
	return paths
}

func TestValueGet(t *testing.T) {
	var gauges = metrics.MemStats{}
	var paths = make([]string, 0)
	var readPath = make([]string, 0)

	gaugeVal := reflect.ValueOf(gauges)
	for i := 0; i < gaugeVal.NumField(); i++ {
		url := fmt.Sprintf("%s/%s/%v/%v", path, gauges, gaugeVal.Field(i).Type().Name(), gaugeVal.Field(i))
		paths = append(paths, url)
		readPath = append(readPath, fmt.Sprintf("/value/%s/%s", gauges, gaugeVal.Field(i).Type().Name()))
	}
	expectSucces := response{200, "text/plain", ""}

	correctPaths := CorrectPaths()

	for i := 0; i < len(paths); i++ {
		t.Run(paths[i], func(t *testing.T) {
			//Running server and executing request
			res := runPost(paths[i])
			//Defering to close body
			defer res.Body.Close()

			//Check response code
			tests.AssertEqualValues(t, expectSucces.code, res.StatusCode)

			//Trying to read body
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			//Check for Content-type values
			tests.AssertHeader(t, res, "Content-Type", expectSucces.contentType)
			//Check for body response
			tests.AssertEqualValues(t, expectSucces.response, string(resBody))

			res = runGet(correctPaths[i])
			//Defering to close body
			defer res.Body.Close()

			//Trying to read body
			resBody1, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("%s\n\n%s\n\n\n\n", resBody1, paths[i])
		})
	}
}
