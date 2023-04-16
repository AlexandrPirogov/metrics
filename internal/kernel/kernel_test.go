package kernel

import (
	"errors"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/memtrack/metrics"
	"testing"

	"github.com/stretchr/testify/assert"
)

type storeStub struct {
	Tuples map[string]map[string]tuples.Tupler
}

func (s *storeStub) Write(t tuples.Tupler) (tuples.Tupler, error) {
	mtype := tuples.ExtractString("type", t)
	switch mtype {
	case "gauge":
		return t, nil
	case "counter":
		return tuples.NewTuple(), nil
	default:
		return tuples.NewTuple(), errors.New("type not exists")
	}
}

func (s *storeStub) Read(t tuples.Tupler) ([]tuples.Tupler, error) {
	return []tuples.Tupler{tuples.NewTuple()}, nil
}

func NewStub() Storer {
	s := &storeStub{}
	return s
}

func TestWriteGaugeTuple(t *testing.T) {
	val := metrics.NumGC(123.123123123)
	sut := tuples.Tuple{
		Fields: map[string]interface{}{
			"name": "qwe",
			"type": "gauge",
			"val":  &val,
		},
	}
	stub := &storeStub{
		make(map[string]map[string]tuples.Tupler),
	}

	actual, err := Write(stub, sut)

	assert.Nil(t, err)
	assert.EqualValues(t, sut, actual)

}

/*
func TestReadGaugeTuple(t *testing.T) {
	val := metrics.NumGC(123.123123123)
	sut := tuples.Tuple{
		Fields: map[string]interface{}{
			"name": "qwe",
			"type": "gauge",
			"val":  &val,
		},
	}
	stub := storeStub{
		make(map[string]map[string]tuples.Tupler),
	}
	stub.Tuples["gauge"] = make(map[string]tuples.Tupler)
	stub.Tuples["gauge"]["qwe"] = sut

	actuals, err := Read(&stub, sut)

	for _, actual := range actuals {

		assert.Nil(t, err)
		assert.EqualValues(t, sut, actual)

	}
}*/
