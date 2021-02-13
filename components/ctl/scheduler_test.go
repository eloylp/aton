// +build unit

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

func TestHeapBasedDetectorPriorityQueue(t *testing.T) {
	q := ctl.NewHeapDetectorPriorityQueue()
	t.Run("Assert order of elements", AssertQueueOrdering(q))
	q = ctl.NewHeapDetectorPriorityQueue()
	t.Run("Assert update of existing element with score updated", AssertQueueUpdate(q))
	q = ctl.NewHeapDetectorPriorityQueue()
	t.Run("Assert remove of existing element", AssertQueueRemove(q))
	q = ctl.NewHeapDetectorPriorityQueue()
	t.Run("Assert queu returns nil with no elements", AssertQueueEmptyNil(q))
}

func AssertQueueOrdering(q ctl.DetectorPriorityQueue) func(t *testing.T) {
	return func(t *testing.T) {
		q.Upsert(AverageUtilizedDetector())
		q.Upsert(FullNetworkUtilizedDetector())
		assert.Equal(t, AverageUtilized, q.Next().UUID)
	}
}

func AssertQueueUpdate(q ctl.DetectorPriorityQueue) func(t *testing.T) {
	return func(t *testing.T) {
		d0 := AverageUtilizedDetector()
		d1 := FullNetworkUtilizedDetector()
		d2 := MidUtilizedDetector()

		q.Upsert(d0)
		q.Upsert(d1)
		q.Upsert(d2)

		d3 := LeastUtilizedDetector()
		d3.UUID = FullNetworkUtilized
		q.Upsert(d3)

		assert.Equal(t, 3, q.Len())
		assert.Equal(t, d3, q.Next())
	}
}

func AssertQueueRemove(q ctl.DetectorPriorityQueue) func(t *testing.T) {
	return func(t *testing.T) {
		d0 := AverageUtilizedDetector()
		d1 := FullNetworkUtilizedDetector()
		d2 := MidUtilizedDetector()

		q.Upsert(d0)
		q.Upsert(d1)
		q.Upsert(d2)

		err := q.Remove(d2.UUID)
		assert.NoError(t, err)
		assert.Equal(t, 2, q.Len())
		assert.Equal(t, d0, q.Next())
	}
}

func AssertQueueEmptyNil(q ctl.DetectorPriorityQueue) func(t *testing.T) {
	return func(t *testing.T) {
		assert.Equal(t, 0, q.Len())
		assert.Nil(t, q.Next())
	}
}
