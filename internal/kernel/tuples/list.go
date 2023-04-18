package tuples

import "encoding/json"

type TupleList struct {
	tuples []Tupler
}

func (t TupleList) Head() Tupler {
	if len(t.tuples) == 0 {
		return nil
	}
	return t.tuples[0]
}

func (t TupleList) Tail() TupleList {
	if !t.Next() {
		return TupleList{}
	}
	t.tuples = t.tuples[1:]
	return t
}

func (t TupleList) Next() bool {
	return len(t.tuples) > 0
}

func (t TupleList) Add(elem Tupler) TupleList {
	if t.tuples == nil {
		t.tuples = make([]Tupler, 0)
	}
	t.tuples = append(t.tuples, elem)
	return t
}

func (t TupleList) HeadTail() (Tupler, TupleList) {
	return t.Head(), t.Tail()
}

func (t TupleList) Merge(toMerge TupleList) TupleList {
	t.tuples = append(t.tuples, toMerge.tuples...)
	return t
}

func MarshalTupleList(tail TupleList, acc []byte) []byte {
	if !tail.Next() {
		return acc
	}
	head, tail := tail.HeadTail()
	bytes, _ := json.Marshal(head)
	return append(bytes, MarshalTupleList(tail, acc)...)
}

func Iterate(list TupleList, fun interface{}) {
	if !list.Next() {

	}
}
