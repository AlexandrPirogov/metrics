package tuples

func NewTuple() Tuple {
	return Tuple{
		Fields: make(map[string]interface{}),
	}
}

type Tupler interface {
	ToTuple() Tupler
	SetField(key string, value interface{}) Tupler
	GetField(key string) (interface{}, bool)
	Aggregate(with Tupler) (Tupler, error)
}

type Tuple struct {
	Fields map[string]interface{}
}

func (t Tuple) SetField(key string, value interface{}) Tupler {
	t.Fields[key] = value
	return t
}

func (t Tuple) GetField(key string) (interface{}, bool) {
	if val, ok := t.Fields[key]; ok {
		return val, true
	}
	return nil, false
}

func (t Tuple) ToTuple() Tupler {
	return t
}

func (t Tuple) Aggregate(with Tupler) (Tupler, error) {
	return t, nil
}
