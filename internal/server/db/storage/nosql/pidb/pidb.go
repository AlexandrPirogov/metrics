package pidb

import (
	"errors"
	"fmt"
	"log"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/metrics"
	"sync"
)

// Local db to imitate storage of metrics
type MemStorage struct {
	Mutex   sync.RWMutex
	Metrics map[string]map[string]tuples.Tupler
}

// Write writes given tuple to Database
//
// Pre-cond: given tuple to write
//
// Post-cond: depends on sucsess
// If success then state was written to database and returned written tuple and error = nil
// Otherwise returns nil and error
func (p *MemStorage) Write(states tuples.TupleList) (tuples.TupleList, error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	return p.writeRec(states, tuples.TupleList{})
}

func (p *MemStorage) writeRec(tail tuples.TupleList, acc tuples.TupleList) (tuples.TupleList, error) {
	if !tail.Next() {
		return acc, nil
	}

	head, tail := tail.HeadTail()
	written, err := p.writeState(head)
	if err != nil {
		return tuples.TupleList{}, err
	}

	return p.writeRec(tail, acc.Add(written))
}

func (p *MemStorage) writeState(state tuples.Tupler) (tuples.Tupler, error) {

	mtype := tuples.ExtractString("type", state)
	mname := tuples.ExtractString("name", state)

	if metrics.IsMetricCorrect(mtype, mname) != nil {
		errMsg := fmt.Sprintf("given not existing metric %s %s\n", mtype, mname)
		return nil, errors.New(errMsg)
	}

	current := p.Metrics[mtype][mname]
	newState, err := state.Aggregate(current)

	if err != nil {
		return nil, err
	}

	p.Metrics[mtype][mname] = newState

	return newState, nil
}

// Read reads tuples from database by given query
//
// Pre-cond: given query tuple
// Post-cond: return tuples that satisfies given query
func (p *MemStorage) Read(state tuples.Tupler) (tuples.TupleList, error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	mtype := tuples.ExtractString("type", state)
	mname := tuples.ExtractString("name", state)
	res := tuples.TupleList{}
	// #TODO and To Mentor: Here will be another refactor later :)
	switch mname {
	case "*":
		// ReadAlllTypes
		switch mtype {
		//Read All Names
		case "*":
			for _, ntype := range p.Metrics {
				for _, metric := range ntype {
					log.Printf("Len of db: %v", metric)
					res = res.Add(metric.ToTuple())
				}
			}
			return res, nil
		default:
			for _, metric := range p.Metrics[mtype] {
				res = res.Add(metric.ToTuple())
			}
			return res, nil
		}
	default:
		if toAppend, ok := p.Metrics[mtype][mname]; ok {
			res = res.Add(toAppend.ToTuple())
			return res, nil
		}
		return tuples.TupleList{}, nil
	}
}

func (p *MemStorage) Ping() error {
	return nil
}
