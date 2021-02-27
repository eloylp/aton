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
		Node     *ctl.Node
		Expected float64
	}{
		{
			Node:     LeastUtilizedNode(),
			Expected: 0.1,
		},
		{
			Node:     OneThirdUtilizedNode(),
			Expected: 0.33,
		},
		{
			Node:     MidUtilizedNode(),
			Expected: 0.5,
		},
		{
			Node:     AverageUtilizedNode(),
			Expected: 0.75,
		},
		{
			Node:     FullUtilizedNode(),
			Expected: math.Inf(1),
		},
		{
			Node:     FullCPUUtilizedNode(),
			Expected: math.Inf(1),
		},
		{
			Node:     FullMemoryUtilizedNode(),
			Expected: math.Inf(1),
		},
		{
			Node:     FullNetworkUtilizedNode(),
			Expected: math.Inf(1),
		},
	}
	for _, tt := range cases {
		t.Run(tt.Node.UUID, func(t *testing.T) {
			result := ctl.NodeUtilizationIndex(tt.Node.Status)
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
func TestScoreNode(t *testing.T) {
	cases := []struct {
		Node     *ctl.Node
		Expected float64
	}{
		{
			Node:     AverageUtilizedNode(),
			Expected: -0.75,
		},
		{
			Node:     FullUtilizedNode(),
			Expected: math.Inf(-1),
		},
	}
	for _, tt := range cases {
		t.Run(tt.Node.UUID, func(t *testing.T) {
			ctl.ScoreNode(tt.Node)
			assert.InDelta(
				t,
				tt.Expected,
				tt.Node.Score,
				0.05,
			)
		})
	}
}

func TestHeapBasedNodePriorityQueue(t *testing.T) {
	q := ctl.NewHeapNodePriorityQueue()
	t.Run("Assert order of elements", AssertQueueOrdering(q))
	q = ctl.NewHeapNodePriorityQueue()
	t.Run("Assert update of existing element with score updated", AssertQueueUpdate(q))
	q = ctl.NewHeapNodePriorityQueue()
	t.Run("Assert remove of existing element", AssertQueueRemove(q))
	q = ctl.NewHeapNodePriorityQueue()
	t.Run("Assert queu returns nil with no elements", AssertQueueEmptyNil(q))
}

func AssertQueueOrdering(q ctl.NodePriorityQueue) func(t *testing.T) {
	return func(t *testing.T) {
		q.Upsert(AverageUtilizedNode())
		q.Upsert(FullNetworkUtilizedNode())
		assert.Equal(t, AverageUtilized, q.Next().UUID)
	}
}

func AssertQueueUpdate(q ctl.NodePriorityQueue) func(t *testing.T) {
	return func(t *testing.T) {
		d0 := AverageUtilizedNode()
		d1 := FullNetworkUtilizedNode()
		d2 := MidUtilizedNode()

		q.Upsert(d0)
		q.Upsert(d1)
		q.Upsert(d2)

		d3 := LeastUtilizedNode()
		d3.UUID = FullNetworkUtilized
		q.Upsert(d3)

		assert.Equal(t, 3, q.Len())
		assert.Equal(t, d3, q.Next())
	}
}

func AssertQueueRemove(q ctl.NodePriorityQueue) func(t *testing.T) {
	return func(t *testing.T) {
		d0 := AverageUtilizedNode()
		d1 := FullNetworkUtilizedNode()
		d2 := MidUtilizedNode()

		q.Upsert(d0)
		q.Upsert(d1)
		q.Upsert(d2)

		err := q.Remove(d2.UUID)
		assert.NoError(t, err)
		assert.Equal(t, 2, q.Len())
		assert.Equal(t, d0, q.Next())
	}
}

func AssertQueueEmptyNil(q ctl.NodePriorityQueue) func(t *testing.T) {
	return func(t *testing.T) {
		assert.Equal(t, 0, q.Len())
		assert.Nil(t, q.Next())
	}
}
