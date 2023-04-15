// Package tuples represents a tuple like in FP languages
package tuples

// NewTuple returns new Tuple
//
// Pre-cond:
//
// Post-cond: returns empty tuple
func NewTuple() Tuple {
	return Tuple{
		Fields: make(map[string]interface{}),
	}
}

type Tupler interface {
	// ToTuple converts type that implements Tupler interface to Tuple
	ToTuple() Tupler

	// SetField adds k/v pair to Tuple.
	//
	// Pre-cond: given key to be set and value for key
	//
	// Post-cond: return new tuple with updated field value
	SetField(key string, value interface{}) Tupler

	// GetField returns value by key of tuple.
	//
	// Pre-cond: given key
	//
	// Post-cond: returns value of field.
	// If field is exists with key returns val and bool = true
	// Otherwise return nil and bool = false
	GetField(key string) (interface{}, bool)

	// Aggregate aggregates to tuples to union them
	//
	// Pre-cond: given tuple to aggregate with
	//
	// Post-cond: returns new Tupler and error
	// If success return Tuplers and error = nil
	// Otherwise return nil and error
	Aggregate(with Tupler) (Tupler, error)
}

// Tuple -- representation of tuple in FP languages
// Just a map where keys are astring and values can be anything
type Tuple struct {
	Fields map[string]interface{}
}

// SetField adds k/v pair to Tuple.
//
// Pre-cond: given key to be set and value for key
//
// Post-cond: return new tuple with updated field value
func (t Tuple) SetField(key string, value interface{}) Tupler {
	t.Fields[key] = value
	return t
}

// GetField returns value by key of tuple.
//
// Pre-cond: given key
//
// Post-cond: returns value of field.
// If field is exists with key returns val and bool = true
// Otherwise return nil and bool = false
func (t Tuple) GetField(key string) (interface{}, bool) {
	if val, ok := t.Fields[key]; ok {
		return val, true
	}
	return nil, false
}

// ToTuple converts type that implements Tupler interface to Tuple
func (t Tuple) ToTuple() Tupler {
	return t
}

// Aggregate aggregates to tuples to union them
//
// Pre-cond: given tuple to aggregate with
//
// Post-cond: returns new Tupler and error
// If success return Tuplers and error = nil
// Otherwise return nil and error
func (t Tuple) Aggregate(with Tupler) (Tupler, error) {
	return t, nil
}
