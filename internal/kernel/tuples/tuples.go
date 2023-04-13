package tuples

type Tupler interface {
	ToTuple() Tuple
}

type Tuple struct {
	fields map[string]interface{}
}

func (t *Tuple) SetField(key string, value interface{}) {
	t.fields[key] = value
}

func (t *Tuple) GetField(key string) interface{} {
	if val, ok := t.fields[key]; ok {
		return val
	}
	return nil
}
