package agent

type RuntimeMetrics struct {
	Alloc         uint64  `json:"alloc,omitempty"`
	BuckHashSys   uint64  `json:"buckHashSys,omitempty"`
	Frees         uint64  `json:"frees,omitempty"`
	GCCPUFraction float64 `json:"GCCPUFraction,omitempty"`
	GCSys         uint64  `json:"GCSys,omitempty"`
	HeapAlloc     uint64  `json:"heapAlloc,omitempty"`
	HeapIdle      uint64  `json:"heapIdle,omitempty"`
	HeapInuse     uint64  `json:"heapInuse,omitempty"`
	HeapObjects   uint64  `json:"heapObjects,omitempty"`
	HeapReleased  uint64  `json:"heapReleased,omitempty"`
	HeapSys       uint64  `json:"heapSys,omitempty"`
	LastGC        uint64  `json:"lastGC,omitempty"`
	Lookups       uint64  `json:"lookups,omitempty"`
	MCacheInuse   uint64  `json:"MCacheInuse,omitempty"`
	MCacheSys     uint64  `json:"MCacheSys,omitempty"`
	MSpanInuse    uint64  `json:"MSpanInuse,omitempty"`
	MSpanSys      uint64  `json:"MSpanSys,omitempty"`
	Mallocs       uint64  `json:"mallocs,omitempty"`
	NextGC        uint64  `json:"nextGC,omitempty"`
	NumForcedGC   uint32  `json:"numForcedGC,omitempty"`
	NumGC         uint32  `json:"numGC,omitempty"`
	OtherSys      uint64  `json:"otherSys,omitempty"`
	PauseTotalNs  uint64  `json:"pauseTotalNs,omitempty"`
	StackInuse    uint64  `json:"stackInuse,omitempty"`
	StackSys      uint64  `json:"stackSys,omitempty"`
	Sys           uint64  `json:"sys,omitempty"`
	TotalAlloc    uint64  `json:"totalAlloc,omitempty"`
}
