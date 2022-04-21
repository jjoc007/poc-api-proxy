package model

type Quota struct {
	SourceIP   string
	TargetPath string
	LimitCalls uint64
}
