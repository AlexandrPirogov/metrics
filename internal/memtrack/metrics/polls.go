package metrics

import "reflect"

type counter int64
type PollCount counter // Count of times when Metrics was collected

type Polls struct {
	PollCount PollCount
}

// Increasing PollCount by 1
// WARNING: PollCount is not reseting
// Overflow may appear
func (p *Polls) Read() int {
	oldP := p.PollCount
	p.PollCount++
	if p.PollCount > oldP {
		return 0
	} else {
		//Overflow appears
		return -1
	}
}

func (p Polls) String() string {
	var tmp counter
	return reflect.TypeOf(tmp).Name()
}
