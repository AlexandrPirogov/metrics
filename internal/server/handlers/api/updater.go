package api

import (
	"log"
	"memtracker/internal/kernel"
	"memtracker/internal/kernel/tuples"
)

// processUpdate updates metric value depends on metric's type
//
// Pre-cond: given metric
//
// Post-cond: return result or processing metric.
// If success, returns slice of bytes and http status = 200
// otherwise returns empty bite slice and corresponging http status
func (d *DefaultHandler) processUpdate(tupls []tuples.Tupler) ([]tuples.Tupler, error) {
	res, err := kernel.Write(d.DB.Storage, tupls)
	if err != nil {
		log.Printf("err while write counter %v", err)
		return []tuples.Tupler{}, err
	}

	return res, nil
}
