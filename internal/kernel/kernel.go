// Package handles job with MetricsState
// It is wrote with declarative paradigm
package kernel

import (
	"memtracker/internal/kernel/tuples"
)

type Replicator interface {
	Write(data []byte)
}

type Storer interface {
	Write(tuple []tuples.Tupler) ([]tuples.Tupler, error)
	Read(cond tuples.Tupler) ([]tuples.Tupler, error)
}

// Write writes tupler to DB and return written state
//
// Pre-cond: given tuples to Write in given Storer
//
// Post-cond: if was written successfully returns NewTuple state and error nil
// Otherwise returns emtyTuple and error
func Write(s Storer, states []tuples.Tupler) ([]tuples.Tupler, error) {
	newStates, err := s.Write(states)
	if err != nil {
		return nil, err
	}

	return newStates, nil
}

func Read(s Storer, state tuples.Tupler) ([]tuples.Tupler, error) {
	states, err := s.Read(state)
	if err != nil {
		return nil, err
	}

	return states, nil
}
