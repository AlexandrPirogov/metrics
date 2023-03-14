package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//TODO ??? move that to better place

// Asserts that two values are equals
func AssertEqualValues(t *testing.T, expected interface{}, actual interface{}) {
	if !assert.EqualValues(t, expected, actual, "Two values must be equals") {
		t.Errorf("values are not equal! Expected %s got %s\n", expected, actual)
	}
}

// Assert headers values for given requests
func AssertHeader(t *testing.T, r *http.Response, header string, expectVal string) {
	AssertEqualValues(t, r.Header.Get(header), expectVal)
}

// Asserts that two values are NOT equals
func AssertNotEqualValues(t *testing.T, expected interface{}, actual interface{}) {
	if assert.EqualValues(t, expected, actual, "Two values must be equals") {
		t.Errorf("we've got unexpected Content-Type! Expected %s got %s\n", expected, actual)
	}
}

// Assert headers values for given requests
func AssertNotHeader(t *testing.T, r *http.Response, header string, expectVal string) {
	AssertNotEqualValues(t, r.Header.Get(header), expectVal)
}

// Require that two values are equals
func RequireEqualValues(t *testing.T, expected interface{}, actual interface{}) {
	require.EqualValues(t, expected, actual, "Two values must be equals")
}

// Require headers values for given requests
func RequireHeader(t *testing.T, r *http.Response, header string, expectVal string) {
	RequireEqualValues(t, r.Header.Get(header), expectVal)
}
