// Package kernel handles job with MetricsState
// It is wrote with declarative paradigm
package kernel

import (
	"log"
	"memtracker/internal/kernel/tuples"
	"time"
)

type Replicator interface {
	Write(data []byte)
}

type Storer interface {
	Write(tuple tuples.TupleList) (tuples.TupleList, error)
	Read(cond tuples.Tupler) (tuples.TupleList, error)
	Ping() error
}

// Write writes tupler to DB and return written state
//
// Pre-cond: given tuples to Write in given Storer
//
// Post-cond: if was written successfully returns NewTuple state and error nil
// Otherwise returns emtyTuple and error
func Write(s Storer, states tuples.TupleList) (tuples.TupleList, error) {
	log.Printf("%s Writing %v", time.Now().Format("15:04:05"), states)
	newStates, err := s.Write(states)
	log.Printf("%s Wrote %v, %v", time.Now().Format("15:04:05"), states, err)
	if err != nil {
		return tuples.TupleList{}, err
	}

	return newStates, nil
}

func Read(s Storer, state tuples.Tupler) (tuples.TupleList, error) {
	log.Printf("%s Reading %v", time.Now().Format("15:04:05"), state)
	states, err := s.Read(state)
	log.Printf("%s Read %v, %v", time.Now().Format("15:04:05"), states, err)
	if err != nil {
		return tuples.TupleList{}, err
	}

	return states, nil
}

// Ping checks is Storer is alive. Health check.
//
// Pre-cond: given Stoerer
//
// Post-cond: makes HealthCheck for given Storer.
// If alive -- return nil, otherwise returns error
func Ping(s Storer) error {
	return s.Ping()
}
