package pidb

import (
	"errors"
	"fmt"
	"log"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/memtrack/metrics"
	"sync"
)

// Local db to imitate storage of metrics
type MemStorage struct {
	Mutex   sync.RWMutex
	Metrics map[string]map[string]tuples.Tupler
}

func (p *MemStorage) Write(state tuples.Tupler) (tuples.Tupler, error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	mtype := p.extractString("type", state)
	mname := p.extractString("name", state)

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

func (p *MemStorage) Read(state tuples.Tupler) ([]tuples.Tupler, error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	//TODO

	mtype := p.extractString("type", state)
	mname := p.extractString("name", state)
	res := make([]tuples.Tupler, 0)
	switch mname {
	case "*":
		// ReadAlllTypes
		switch mtype {
		//Read All Names
		case "*":
			for _, ntype := range p.Metrics {
				for _, metric := range ntype {
					log.Printf("Len of db: %v", metric)
					res = append(res, metric.ToTuple())
				}
			}
			return res, nil
		default:
			for _, metric := range p.Metrics[mtype] {
				res = append(res, metric.ToTuple())
			}
			return res, nil
		}
	default:
		if toAppend, ok := p.Metrics[mtype][mname]; ok {
			res = append(res, toAppend)
			return res, nil
		}
		return []tuples.Tupler{}, errors.New("not found")
	}
}
