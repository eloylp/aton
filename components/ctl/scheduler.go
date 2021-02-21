package ctl

import (
	"container/heap"
	"fmt"
	"math"
)

// DetectorUtilizationIndex calculates a utilization index based on the
// internal consumption of system resources in a detector. The higher
// the value the more busy will be the detector, so less eligible.
// This function will return an infinite positive value in case of
// any resource exceeds the limit threshold, indicating that the detector
// should not be chosen. We cannot assume that this function will always
// return a percentage. Its just an index, the lower, the best.
//
// The current implementation is a simplistic one (see comments) and works
// on a best effort.
func DetectorUtilizationIndex(s *Status) float64 {
	// Find percentages of each part of the system
	loadPercent := s.System.LoadAverage.Avg1 / float64(s.System.CPUCount)
	memPercent := float64(s.System.Memory.UsedMemoryBytes) / float64(s.System.Memory.TotalMemoryBytes)
	networkPercent := float64(s.System.Network.RxBytesSec) / 100e6 // Assumes Gigabit ethernet.

	// Prevent edge cases when any resource its at its max utilization.
	useMaxThreshold := 0.98 // 98% percent of use.
	if loadPercent >= useMaxThreshold || memPercent >= useMaxThreshold || networkPercent >= useMaxThreshold {
		return math.Inf(1)
	}

	// Calculate the average of the utilization percentages of the system.
	return (loadPercent + memPercent + networkPercent) / 3
}

// ScoreDetector calculates a general score for the Detector passed
// and sets the value in the Detector struct.
// This function tends to include all the available calculations
// in order to nurture a priority queue. The more the score the more
// eligible will be the detector. Negative scoring is possible.
func ScoreDetector(d *Detector) {
	d.Score = DetectorUtilizationIndex(d.Status) * -1 // Negative score, as this is utilization.
}

// DetectorPriorityQueue defines the interfaces needed for
// interacting with the scheduler. Multiple implementations based
// on different criteria are expected.
type DetectorPriorityQueue interface {
	// Upsert will add the *Detector passed as argument to the queue.
	// If the element already exists it will replace it.
	Upsert(*Detector)
	Len() int
	Remove(string) error
	// Next must return the next most suitable *Detector for doing some task.
	// When the queue is empty, nil should be returned.
	Next() *Detector
}

// HeapDetectorPriorityQueue is an implementation of DetectorPriorityQueue based on a
// heap. This necessary implements heap.Interface, as we are using the out of the box
// heap of the std lib. Such methods should be used only internally by the std lib.
type HeapDetectorPriorityQueue struct {
	list []*Detector
	uuid map[string]*Detector
}

func NewHeapDetectorPriorityQueue() *HeapDetectorPriorityQueue {
	return &HeapDetectorPriorityQueue{
		uuid: make(map[string]*Detector),
	}
}

func (h *HeapDetectorPriorityQueue) Upsert(detector *Detector) {
	ScoreDetector(detector)
	if _, ok := h.uuid[detector.UUID]; ok {
		if err := h.Remove(detector.UUID); err != nil {
			panic(err)
		}
	}
	h.uuid[detector.UUID] = detector
	heap.Push(h, detector)
}

func (h *HeapDetectorPriorityQueue) Remove(uuid string) error {
	sd, ok := h.uuid[uuid]
	if !ok {
		return fmt.Errorf("scheduler: heap: cannot find uuid %s", uuid)
	}
	heap.Remove(h, sd.Index)
	delete(h.uuid, uuid)
	return nil
}

func (h *HeapDetectorPriorityQueue) Next() *Detector {
	if len(h.list) == 0 {
		return nil
	}
	return heap.Pop(h).(*Detector)
}

func (h *HeapDetectorPriorityQueue) Len() int { return len(h.list) }
func (h *HeapDetectorPriorityQueue) Less(i, j int) bool {
	return h.list[i].Score > h.list[j].Score
}

func (h *HeapDetectorPriorityQueue) Swap(i, j int) {
	h.list[i], h.list[j] = h.list[j], h.list[i]
	h.list[i].Index = i
	h.list[j].Index = j
}

func (h *HeapDetectorPriorityQueue) Push(x interface{}) {
	n := len(h.list)
	d := x.(*Detector)
	d.Index = n
	h.list = append(h.list, d)
}

func (h *HeapDetectorPriorityQueue) Pop() interface{} {
	old := h.list
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	h.list = old[0 : n-1]
	return item
}
