// Package metrics contains all metrics type to collect
package metrics

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type TotalMemory gauge

type FreeMemory gauge

type CPUutilization1 gauge

type OpsUtil struct {
	FreeMemory      FreeMemory
	CPUutilization1 CPUutilization1
	TotalMemory     TotalMemory
}

// Read Increasing PollCount by 1
// WARNING: PollCount is not reseting
// Overflow may appear
func (o *OpsUtil) Read() error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	cp, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}

	o.FreeMemory = FreeMemory(v.Free)
	o.TotalMemory = TotalMemory(v.Total)
	o.CPUutilization1 = CPUutilization1(cp[0])
	return nil
}

func (o OpsUtil) AsMap() map[string]interface{} {
	metrics := make(map[string]interface{}, 3)
	metrics["CPUutilization1"] = float64(o.CPUutilization1)
	metrics["FreeMemory"] = float64(o.FreeMemory)
	metrics["TotalMemory"] = float64(o.TotalMemory)
	return metrics
}

func (o OpsUtil) String() string {
	return "gauge"
}
