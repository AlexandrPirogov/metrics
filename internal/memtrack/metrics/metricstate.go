package metrics

import (
	"errors"
	"log"
	"memtracker/internal/kernel/tuples"
	"strconv"
)

func CreateState(mname, mtype, mvalue string) (tuples.Tupler, error) {
	switch mtype {
	case "gauge":
		return createGaugeState(mname, mtype, mvalue)
	case "counter":
		return createCounterState(mname, mtype, mvalue)
	default:
		return nil, errors.New("given not existing metric's type")
	}
}

func createGaugeState(mname, mtype, mvalue string) (tuples.Tupler, error) {
	float64Val, err := strconv.ParseFloat(mvalue, 64)
	if err != nil {
		return GaugeState{
			Name:  mname,
			Type:  mtype,
			Value: nil,
		}, nil
	}

	return GaugeState{
		Name:  mname,
		Type:  mtype,
		Value: &float64Val,
	}, nil
}

func createCounterState(mname, mtype, mvalue string) (tuples.Tupler, error) {
	int64Val, err := strconv.ParseInt(mvalue, 10, 64)
	if err != nil {
		return CounterState{
			Name:  mname,
			Type:  mtype,
			Value: nil,
		}, err
	}

	return CounterState{
		Name:  mname,
		Type:  mtype,
		Value: &int64Val,
	}, nil
}

// Struct to represent Counter metrics as tuple in FP style
type CounterState struct {
	Name  string `json:"id"`              //Metric name
	Type  string `json:"type"`            // Metric type: gauge or counter
	Value *int64 `json:"delta,omitempty"` //Metric's val if passing counter
}

// ToTuple convers CounterState to tuple.Tuple
//
// Pre-cond:
//
// Post-cond: tuples instance are returned
func (c CounterState) ToTuple() tuples.Tupler {
	tuple := CounterState{
		Name:  c.Name,
		Type:  c.Type,
		Value: c.Value,
	}
	return tuple
}

// Changes property of CounterState for given value
//
// Pre-cond: given existing key and appropriate value
//
// Post-cond: returns new CounterState, otherwise throws an error
func (c CounterState) SetField(key string, value interface{}) tuples.Tupler {
	switch key {
	case "name":
		c.Name = value.(string)
	case "type":
		c.Type = value.(string)
	case "value":
		c.Value = value.(*int64)
	default:
		log.Fatalf("given not existing field %s for counter state!", key)
	}
	return c
}

// Read property of CounterState for given value
//
// Pre-cond: given existing key
//
// Post-cond: returns corresponding value and true, otherwise empty string and false
func (c CounterState) GetField(key string) (interface{}, bool) {
	switch key {
	case "name":
		return c.Name, true
	case "type":
		return c.Type, true
	case "value":
		if c.Value == nil {
			return nil, false
		}
		return c.Value, true
	default:
		return "", false
	}
}

// Aggregate doing nothing for gaugeState
//
// Pre-cond: given tupler to aggregate with
//
// Post-cond: if given TODO
func (c CounterState) Aggregate(with tuples.Tupler) (tuples.Tupler, error) {
	val, ok := c.GetField("value")
	if !ok || val == nil {
		return nil, errors.New("value must exists for writing metric")
	}

	if with == nil {
		return c, nil
	}

	val, ok = with.GetField("value")
	if !ok || val == nil {
		return c, nil
	}

	val64 := val.(*int64)
	newVal := *(c.Value) + *(val64)
	c.Value = &newVal
	return c, nil
}

// Struct to represent gauge metrics as tuple in FP style
type GaugeState struct {
	Name  string   `json:"id"`              //Metric name
	Type  string   `json:"type"`            // Metric type: gauge or counter
	Value *float64 `json:"value,omitempty"` //Metric's val if passing counter
}

// ToTuple convers GaugeState to tuple.Tuple
//
// Pre-cond:
//
// Post-cond: tuples instance are returned
func (g GaugeState) ToTuple() tuples.Tupler {
	tuple := GaugeState{
		Name:  g.Name,
		Type:  g.Type,
		Value: g.Value,
	}
	return tuple
}

// Changes property of CounterState for given value
//
// Pre-cond: given existing key and appropriate value
//
// Post-cond: returns new CounterState, otherwise throws an error
func (g GaugeState) SetField(key string, value interface{}) tuples.Tupler {
	switch key {
	case "name":
		g.Name = value.(string)
	case "type":
		g.Type = value.(string)
	case "value":
		g.Value = value.(*float64)
	default:
		log.Fatalf("given not existing field %s for gauge state!", key)
	}
	return g
}

// Read property of CounterState for given value
//
// Pre-cond: given existing key
//
// Post-cond: returns corresponding value and true, otherwise empty string and false
func (g GaugeState) GetField(key string) (interface{}, bool) {
	switch key {
	case "name":
		return g.Name, true
	case "type":
		return g.Type, true
	case "value":
		if g.Value == nil {
			return nil, false
		}
		return g.Value, true
	default:
		return "", false
	}
}

// Aggregate doing nothing for gaugeState
//
// Pre-cond: given tupler to aggregate with
//
// Post-cond: ...
func (g GaugeState) Aggregate(with tuples.Tupler) (tuples.Tupler, error) {
	return g, nil
}
