package model

type CallsStatsSummary struct {
	SourceIP   string
	TargetPath string
	Year       uint8
	Month      uint8
	TotalCalls uint64
}
