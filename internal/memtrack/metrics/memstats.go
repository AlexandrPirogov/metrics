package metrics

import (
	"math/rand"
	"reflect"
	"runtime"
)

type gauge float64

// The amount of memory allocated by the application.
type Alloc gauge

// The amount of memory used by Go's internal hash table implementation
type BuckHashSys gauge

// The amount of memory freed by the application.
type Frees gauge

// The fraction of CPU time used by the garbage collector
type GCCPUFraction gauge

// The amount of memory used by the garbage collector.
type GCSys gauge

// The amount of memory allocated to the heap by the application
type HeapAlloc gauge

// The amount of idle (unused) heap memory.
type HeapIdle gauge

// The amount of heap memory currently in use.
type HeapInuse gauge

// The number of objects currently in the heap.
type HeapObjects gauge

// The amount of heap memory that has been returned to the operating system.
type HeapReleased gauge

// The total size of the heap memory.
type HeapSys gauge

// The time (in nanoseconds) since the last garbage collection.
type LastGC gauge

// The number of pointer lookups performed by the garbage collector
type Lookups gauge

// The amount of memory used by runtime mcache structures that are currently being used
type MCacheInuse gauge

// The total size of the memory reserved for runtime mcache structures
type MCacheSys gauge

// The amount of memory used by runtime mspan structures that are currently being used
type MSpanInuse gauge

// The total size of the memory reserved for runtime mspan structures
type MSpanSys gauge

// The number of heap memory allocations made by the application.
type Mallocs gauge

// The estimated heap size (in bytes) when the next garbage collection will occur.
type NextGC gauge

// The number of garbage collections that have been forced by the application.
type NumForcedGC gauge

// The number of garbage collections that have been performed by the application.
type NumGC gauge

// The amount of memory used by other system-level activities (e.g. network I/O
type OtherSys gauge

// The total time (in nanoseconds) spent by the garbage collector in performing pauses (i.e. when the application is stopped).
type PauseTotalNs gauge

// The amount of stack memory currently in use by the application.
type StackInuse gauge

// The total size of the stack memory.
type StackSys gauge

// The total memory used by the application (including heap, stack, and other system-level activities).
type Sys gauge

// The total amount of memory allocated by the application, including memory that has been freed.
type TotalAlloc gauge

// The random value
type RandomValue gauge

// Container that holds all mem stats metrics
type MemStats struct {
	Alloc         Alloc
	BuckHashSys   BuckHashSys
	Frees         Frees
	GCCPUFraction GCCPUFraction
	GCSys         GCSys
	HeapAlloc     HeapAlloc
	HeapIdle      HeapIdle
	HeapInuse     HeapInuse
	HeapObjects   HeapObjects
	HeapReleased  HeapReleased
	HeapSys       HeapSys
	LastGC        LastGC
	Lookups       Lookups
	MCacheInuse   MCacheInuse
	MCacheSys     MCacheSys
	MSpanInuse    MSpanInuse
	MSpanSys      MSpanSys
	Mallocs       Mallocs
	NextGC        NextGC
	NumForcedGC   NumForcedGC
	NumGC         NumGC
	OtherSys      OtherSys
	PauseTotalNs  PauseTotalNs
	StackInuse    StackInuse
	StackSys      StackSys
	Sys           Sys
	TotalAlloc    TotalAlloc
	RandomValue   RandomValue
}

// Pre-cond:
// Post-cond: metrics updated using package runtime
// Reflect could be used here
func (m *MemStats) Read() int {
	var runtimeMemStat runtime.MemStats
	runtime.ReadMemStats(&runtimeMemStat)
	m.Alloc = Alloc(runtimeMemStat.Alloc)
	m.BuckHashSys = BuckHashSys(runtimeMemStat.BuckHashSys)
	m.Frees = Frees(runtimeMemStat.Frees)
	m.GCCPUFraction = GCCPUFraction(runtimeMemStat.GCCPUFraction)
	m.GCSys = GCSys(runtimeMemStat.GCSys)
	m.HeapAlloc = HeapAlloc(runtimeMemStat.HeapAlloc)
	m.HeapIdle = HeapIdle(runtimeMemStat.HeapIdle)
	m.HeapInuse = HeapInuse(runtimeMemStat.HeapInuse)
	m.HeapObjects = HeapObjects(runtimeMemStat.HeapObjects)
	m.HeapReleased = HeapReleased(runtimeMemStat.HeapReleased)
	m.HeapSys = HeapSys(runtimeMemStat.HeapSys)
	m.LastGC = LastGC(runtimeMemStat.LastGC)
	m.Lookups = Lookups(runtimeMemStat.Lookups)
	m.MCacheInuse = MCacheInuse(runtimeMemStat.MCacheInuse)
	m.MCacheSys = MCacheSys(runtimeMemStat.MCacheSys)
	m.MSpanInuse = MSpanInuse(runtimeMemStat.MSpanInuse)
	m.MSpanSys = MSpanSys(runtimeMemStat.MSpanSys)
	m.Mallocs = Mallocs(runtimeMemStat.Mallocs)
	m.NextGC = NextGC(runtimeMemStat.NextGC)
	m.NumForcedGC = NumForcedGC(runtimeMemStat.NumForcedGC)
	m.NumGC = NumGC(runtimeMemStat.NumGC)
	m.OtherSys = OtherSys(runtimeMemStat.OtherSys)
	m.PauseTotalNs = PauseTotalNs(runtimeMemStat.PauseTotalNs)
	m.StackInuse = StackInuse(runtimeMemStat.StackInuse)
	m.StackSys = StackSys(runtimeMemStat.StackSys)
	m.Sys = Sys(runtimeMemStat.Sys)
	m.TotalAlloc = TotalAlloc(runtimeMemStat.TotalAlloc)
	m.RandomValue = RandomValue(rand.Float64())
	return 0
}

func (m MemStats) String() string {
	var tmp gauge
	return reflect.TypeOf(tmp).Name()
}
