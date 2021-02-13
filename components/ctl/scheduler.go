package ctl

import (
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
