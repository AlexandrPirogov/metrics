package tuples

// ExtractString extracts string field from tuple
//
// Pre-cond: given field to extract and tuple to extract from
//
// Post-cond: extracts string value of fields.
// If field no exists or field is not string, return empty string
func ExtractString(field string, t Tupler) string {
	f, ok := t.GetField(field)
	if !ok {
		return ""
	}

	return f.(string)
}

// ExtractString extracts pointer to int64 field from tuple
//
// Pre-cond: given field to extract and tuple to extract from
//
// Post-cond: extracts pointer to int64  value of fields.
// If field no exists or field is not pointer to int64 , return nil
func ExtractInt64Pointer(field string, t Tupler) *int64 {
	f, ok := t.GetField("value")
	if !ok {
		return nil
	}
	return f.(*int64)
}

// ExtractString extracts pointer to float64 field from tuple
//
// Pre-cond: given field to extract and tuple to extract from
//
// Post-cond: extracts pointer to float64  value of fields.
// If field no exists or field is not pointer to float64 , return nil
func ExtractFloat64Pointer(field string, t Tupler) *float64 {
	f, ok := t.GetField("value")
	if !ok {
		return nil
	}
	return f.(*float64)
}
