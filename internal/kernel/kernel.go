// Package handles job with MetricsState
// It is wrote with declarative paradigm
package kernel

import (
	"log"
	"memtracker/internal/kernel/tuples"
)

type Replicator interface {
	Write(data []byte)
}

type Storer interface {
	Write(tuple tuples.TupleList) (tuples.TupleList, error)
	Read(cond tuples.Tupler) (tuples.TupleList, error)
}

// Write writes tupler to DB and return written state
//
// Pre-cond: given tuples to Write in given Storer
//
// Post-cond: if was written successfully returns NewTuple state and error nil
// Otherwise returns emtyTuple and error
func Write(s Storer, states tuples.TupleList) (tuples.TupleList, error) {
	newStates, err := s.Write(states)

	if err != nil {
		log.Printf("err while writing state %v", err)
		return tuples.TupleList{}, err
	}

	return newStates, nil
}

func Read(s Storer, state tuples.Tupler) (tuples.TupleList, error) {
	states, err := s.Read(state)
	if err != nil {
		log.Printf("err while reading state %v", err)
		return tuples.TupleList{}, err
	}

	return states, nil
}
