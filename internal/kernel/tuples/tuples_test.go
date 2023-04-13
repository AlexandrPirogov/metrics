package tuples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTupleSetKey(t *testing.T) {
	cases := map[string]interface{}{
		"a": 1,
		"b": "test",
		"something": struct {
			Name  string
			Val   float64
			Delta int64
		}{
			"alex",
			1.11111,
			1,
		},
	}
	sut := Tuple{
		Fields: make(map[string]interface{}),
	}

	for k, expectedValue := range cases {
		sut.SetField(k, expectedValue)
		actual := sut.GetField(k)

		assert.EqualValues(t, expectedValue, actual)
	}

}
