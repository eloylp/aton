package ctl

import (
	"time"
)

type Detector struct {
	UUID   string
	Index  int
	Score  float64
	Status *Status
}

type Status struct {
	Description string
	Capturers   []*Capturer
	System      System
}

type System struct {
	CPUCount    int
	Network     Network
	LoadAverage LoadAverage
	Memory      Memory
}

type Network struct {
	RxBytesSec uint64
	TxBytesSec uint64
}

type LoadAverage struct {
	Avg1  float64
	Avg5  float64
	Avg15 float64
}

type Memory struct {
	UsedMemoryBytes  uint64
	TotalMemoryBytes uint64
}

type Capturer struct {
	UUID   string
	URL    string
	Status string
}

type Result struct {
	DetectorUUID  string
	Recognized    []string
	TotalEntities int32
	RecognizedAt  time.Time
	CapturedAt    time.Time
}
