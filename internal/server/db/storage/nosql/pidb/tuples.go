package pidb

import "memtracker/internal/kernel/tuples"

func (p *MemStorage) extractString(field string, t tuples.Tupler) string {
	f, ok := t.GetField(field)
	if !ok {
		return ""
	}

	return f.(string)
}
