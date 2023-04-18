package tuples

import "encoding/json"

// TupleList is a container of instance that dispatches Tupler interface
type TupleList struct {
	tuples []Tupler
}

// Head return first elem of TupleList. If there is no elems return nil
//
// Pre-cond:
//
// Post-cond: if TupleList isn't empty returns first element of list
// Otherwise returns nil
func (t TupleList) Head() Tupler {
	if len(t.tuples) == 0 {
		return nil
	}
	return t.tuples[0]
}

// Head Tail returns new TupleList except of it head
//
// Pre-cond:
//
// Post-cond: returns new TupleList except of head if there are any elements in tail
// Otherwise returns empty TupleList
func (t TupleList) Tail() TupleList {
	if !t.Next() {
		return TupleList{}
	}
	t.tuples = t.tuples[1:]
	return t
}

// Next tells if there are elems in list
//
// Pre-cond:
//
// Post-cond: if TupleList has more than 0 elems returns true
// Otherwise returns false
func (t TupleList) Next() bool {
	return len(t.tuples) > 0
}

// Len return count of elems in TupleList
//
// Pre-cond:
//
// Post-cond: returns count of elems in TupleList
func (t TupleList) Len() int {
	if t.tuples == nil {
		return 0
	}
	return len(t.tuples)
}

// Add adds elem to TupleList
//
// Pre-cond: given Tupler element
//
// Post-cond: elem was added to TupleList. Returns new TupleList with given elem
func (t TupleList) Add(elem Tupler) TupleList {
	if t.tuples == nil {
		t.tuples = make([]Tupler, 0)
	}
	t.tuples = append(t.tuples, elem)
	return t
}

// HeadTail returns head and tail of TupleList
//
// Pre-cond:
//
// Post-cond: returns head and rest TupleList
func (t TupleList) HeadTail() (Tupler, TupleList) {
	return t.Head(), t.Tail()
}

// Merge merges two TupleLists into one
//
// Pre-cond:
//
// Post-cond: returns new TupleList that has elems of merging Tuplelists
func (t TupleList) Merge(toMerge TupleList) TupleList {
	t.tuples = append(t.tuples, toMerge.tuples...)
	return t
}

// AsSlice returns slice of elems in TupleList
//
// Pre-cond:
//
// Post-cond: returns slice of elems in TupleList
func (t TupleList) AsSlice() []Tupler {
	return t.tuples
}

// MarshalTupleList marshals all elems in TupleList
//
// Pre-cond: given TupleList and acc to store result
//
// Post-cond: returns marshaled elems
func MarshalTupleList(tail TupleList, acc []byte) []byte {
	slice := tail.AsSlice()
	if len(slice) == 1 {
		body, _ := json.Marshal(slice[0])
		return body
	}

	body, _ := json.Marshal(slice)
	return body
}
