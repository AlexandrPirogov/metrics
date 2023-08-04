package kernel

import (
	"errors"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/metrics"
	"testing"

	"github.com/stretchr/testify/assert"
)

type storeStub struct {
	Tuples map[string]map[string]tuples.Tupler
}

func (s *storeStub) Write(tupls tuples.TupleList) (tuples.TupleList, error) {
	for tupls.Next() {
		head := tupls.Head()
		mtype := tuples.ExtractString("type", head)
		switch mtype {
		case "gauge":
			return tuples.TupleList{}.Add(head), nil
		case "counter":
			return tuples.TupleList{}.Add(head), nil
		default:
			return tuples.TupleList{}, errors.New("type not exists")
		}
	}
	return tuples.TupleList{}, errors.New("type not exists")
}

func (s *storeStub) Read(t tuples.Tupler) (tuples.TupleList, error) {
	return tuples.TupleList{}, nil
}

func (s *storeStub) Ping() error {
	return nil
}

func NewStub() Storer {
	s := &storeStub{}
	return s
}

func TestWriteGaugeTuple(t *testing.T) {
	val := metrics.NumGC(123.123123123)
	sut := tuples.TupleList{}.Add(
		tuples.Tuple{
			Fields: map[string]interface{}{
				"name": "qwe",
				"type": "gauge",
				"val":  &val,
			},
		})

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
