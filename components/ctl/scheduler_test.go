package ctl_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/components/ctl"
)

// nolint: scopelint
func TestUtilizationIndex(t *testing.T) {
	cases := []struct {
		Detector *ctl.Detector
		Expected float64
	}{
		{
			Detector: LeastUtilizedDetector(),
			Expected: 0.1,
		},
		{
			Detector: OneThirdUtilizedDetector(),
			Expected: 0.33,
		},
		{
			Detector: MidUtilizedDetector(),
			Expected: 0.5,
		},
		{
			Detector: AverageUtilizedDetector(),
			Expected: 0.75,
		},
		{
			Detector: FullUtilizedDetector(),
			Expected: math.Inf(1),
		},
		{
			Detector: FullCPUUtilizedDetector(),
			Expected: math.Inf(1),
		},
		{
			Detector: FullMemoryUtilizedDetector(),
			Expected: math.Inf(1),
		},
		{
			Detector: FullNetworkUtilizedDetector(),
			Expected: math.Inf(1),
		},
	}
	for _, tt := range cases {
		t.Run(tt.Detector.UUID, func(t *testing.T) {
			result := ctl.DetectorUtilizationIndex(tt.Detector.Status)
			assert.InDelta(
				t,
				tt.Expected,
				result,
				0.05,
			)
		})
	}
}

// nolint: scopelint
func TestScoreDetector(t *testing.T) {
	cases := []struct {
		Detector *ctl.Detector
		Expected float64
	}{
		{
			Detector: AverageUtilizedDetector(),
			Expected: -0.75,
		},
		{
			Detector: FullUtilizedDetector(),
			Expected: math.Inf(-1),
		},
	}
	for _, tt := range cases {
		t.Run(tt.Detector.UUID, func(t *testing.T) {
			ctl.ScoreDetector(tt.Detector)
			assert.InDelta(
				t,
				tt.Expected,
				tt.Detector.Score,
				0.05,
			)
		})
	}
}
