package main_test

import (
	"fmt"
	"io"
	"memtracker/internal/memtrack/metrics"
	"memtracker/internal/server/handlers"
	"memtracker/internal/tests"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type response struct {
	code        int
	contentType string
	response    string
}

const host string = "http://localhost:8080"
const path string = "/update"

func TestUpdateHandlerIncorrectPath(t *testing.T) {
	expectFail := response{http.StatusNotFound, "", ""}
	incorrectPaths := []string{
		host + path + "/qwe/asd",
		host + path + "/qwe",
		host + path + "/A/",
		host + path + "/A/",
		host + path + "/C/",
	}

	for _, url := range incorrectPaths {
		t.Run(url, func(t *testing.T) {
			//Running server and executing request
			res := runPost(url).Result()

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

func TestUPdateHandlerCorrectPath(t *testing.T) {
	expectSucces := response{200, "text/plain", ""}

	correctPaths := CorrectPaths()

	for _, url := range correctPaths {
		t.Run(url, func(t *testing.T) {
			//Running server and executing request
			res := runPost(url).Result()

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
		})
	}
}

func runPost(url string) *httptest.ResponseRecorder {
	//Creating request to execute
	request := httptest.NewRequest(http.MethodPost, url, nil)
	request.Header.Set("Content-Type", "text/plain")

	//Creating recorder
	recorder := httptest.NewRecorder()

	//Define handler
	testingHandler := http.HandlerFunc(handlers.UpdateHandler)

	//Running server
	testingHandler.ServeHTTP(recorder, request)
	return recorder
}

func CorrectPaths() []string {
	var gauges = metrics.MemStats{}
	var paths = make([]string, 0)

	gaugeVal := reflect.ValueOf(gauges)
	for i := 0; i < gaugeVal.NumField(); i++ {
		url := host + path + fmt.Sprintf("/%s/%v/%v", gauges, gaugeVal.Field(i).Type().Name(), gaugeVal.Field(i))
		paths = append(paths, url)
	}
	return paths
}
