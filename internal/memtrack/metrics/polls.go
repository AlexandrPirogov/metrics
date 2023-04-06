package metrics

import (
	"fmt"
)

type counter int64
type PollCount counter // Count of times when Metrics was collected

type Polls struct {
	PollCount PollCount
}

// Read Increasing PollCount by 1
// WARNING: PollCount is not reseting
// Overflow may appear
func (p *Polls) Read() error {
	oldP := p.PollCount
	p.PollCount++
	//Checks for overflow
	if p.PollCount > oldP {
		return nil
	} else {
		//Overflow appears
		return fmt.Errorf("overflow appeared")
	}
}

func (p Polls) AsMap() map[string]interface{} {
	metrics := make(map[string]interface{}, 1)
	metrics["PollCount"] = float64(p.PollCount)
	return metrics
}

/*
func (p Polls) AsMap() map[string]interface{} {
	metrics := make(map[string]interface{}, 28)
	metrics["pollcount"] = PollCount(p.PollCount)
	return metrics
}*/

func (p Polls) String() string {
	return "counter"
}
