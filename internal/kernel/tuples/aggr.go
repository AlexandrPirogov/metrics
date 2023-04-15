package tuples

func Append[T int64 | float64](f string, t, t1 Tupler) Tupler {
	//TODO
	f1, _ := t.GetField(f)
	f2, _ := t1.GetField(f)
	t.SetField(f, f1.(T)+f2.(T))
	return t
}
