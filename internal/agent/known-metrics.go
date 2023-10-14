package agent

type RuntimeMetrics struct {
	Alloc         uint64  `json:"Alloc"`
	BuckHashSys   uint64  `json:"BuckHashSys"`
	Frees         uint64  `json:"Frees"`
	GCCPUFraction float64 `json:"GCCPUFraction"`
	GCSys         uint64  `json:"GCSys"`
	HeapAlloc     uint64  `json:"HeapAlloc"`
	HeapIdle      uint64  `json:"HeapIdle"`
	HeapInuse     uint64  `json:"HeapInuse"`
	HeapObjects   uint64  `json:"HeapObjects"`
	HeapReleased  uint64  `json:"HeapReleased"`
	HeapSys       uint64  `json:"HeapSys"`
	LastGC        uint64  `json:"LastGC"`
	Lookups       uint64  `json:"Lookups"`
	MCacheInuse   uint64  `json:"MCacheInuse"`
	MCacheSys     uint64  `json:"MCacheSys"`
	MSpanInuse    uint64  `json:"MSpanInuse"`
	MSpanSys      uint64  `json:"MSpanSys"`
	Mallocs       uint64  `json:"Mallocs"`
	NextGC        uint64  `json:"NextGC"`
	NumForcedGC   uint32  `json:"NumForcedGC"`
	NumGC         uint32  `json:"NumGC"`
	OtherSys      uint64  `json:"OtherSys"`
	PauseTotalNs  uint64  `json:"PauseTotalNs"`
	StackInuse    uint64  `json:"StackInuse"`
	StackSys      uint64  `json:"StackSys"`
	Sys           uint64  `json:"Sys"`
	TotalAlloc    uint64  `json:"TotalAlloc"`
}
