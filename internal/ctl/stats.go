package ctl

import (
	"sync/atomic"
)

type Stats struct {
	currentDetectors int32
	currentCapturers int32
	processedSuccess int64
	processedFailed  int64
}

func (s *Stats) CurrentCapturers() int32 {
	return atomic.LoadInt32(&s.currentCapturers)
}

func (s *Stats) CurrentDetectors() int32 {
	return atomic.LoadInt32(&s.currentDetectors)
}

func (s *Stats) Processed() int64 {
	return atomic.LoadInt64(&s.processedSuccess) + atomic.LoadInt64(&s.processedFailed)
}

func (s *Stats) ProcessedSuccess() int64 {
	return atomic.LoadInt64(&s.processedSuccess)
}

func (s *Stats) ProcessedFailed() int64 {
	return atomic.LoadInt64(&s.processedFailed)
}

func (s *Stats) IncSuccess() {
	atomic.AddInt64(&s.processedSuccess, 1)
}

func (s *Stats) IncFailed() {
	atomic.AddInt64(&s.processedFailed, 1)
}

func (s *Stats) IncCapturers() {
	atomic.AddInt32(&s.currentCapturers, 1)
}

func (s *Stats) IncDetectors(i int32) {
	atomic.AddInt32(&s.currentDetectors, i)
}
