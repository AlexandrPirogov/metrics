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

func (s *storeStub) Write(tupls []tuples.Tupler) ([]tuples.Tupler, error) {
	for _, t := range tupls {
		mtype := tuples.ExtractString("type", t)
		switch mtype {
		case "gauge":
			return []tuples.Tupler{t}, nil
		case "counter":
			return []tuples.Tupler{}, nil
		default:
			return []tuples.Tupler{}, errors.New("type not exists")
		}
	}
	return []tuples.Tupler{}, errors.New("type not exists")
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
	sut := []tuples.Tupler{
		tuples.Tuple{
			Fields: map[string]interface{}{
				"name": "qwe",
				"type": "gauge",
				"val":  &val,
			},
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
