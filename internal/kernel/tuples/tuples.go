package tuples

type Tupler interface {
	ToTuple() Tuple
}

type Tuple struct {
	Fields map[string]interface{}
}

func (t *Tuple) SetField(key string, value interface{}) {
	t.Fields[key] = value
}

func (t *Tuple) GetField(key string) interface{} {
	if val, ok := t.Fields[key]; ok {
		return val
	}
	return nil
}
