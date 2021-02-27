package ctl_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/eloylp/aton/components/ctl"
)

var (
	random  = rand.New(rand.NewSource(time.Now().Unix())) //nolint:gosec
	nodeSet = [...]*ctl.Node{
		LeastUtilizedNode(),
		OneThirdUtilizedNode(),
		MidUtilizedNode(),
		AverageUtilizedNode(),
		FullUtilizedNode(),
	}
)

func Benchmark_HeapPriorityQueue(b *testing.B) {
	q := ctl.NewHeapNodePriorityQueue()
	b.Run("Benchmark push operation with preload of 10 elements", BenchmarkingPush(q, 10))
	q = ctl.NewHeapNodePriorityQueue()
	b.Run("Benchmark push operation with preload of 100 elements", BenchmarkingPush(q, 100))
	q = ctl.NewHeapNodePriorityQueue()
	b.Run("Benchmark push operation with preload of 100000 elements", BenchmarkingPush(q, 100000))

	q = ctl.NewHeapNodePriorityQueue()
	b.Run("Benchmark pull operation with preload of 10 elements", BenchmarkingPull(q, 10))
	q = ctl.NewHeapNodePriorityQueue()
	b.Run("Benchmark pull operation with preload of 100 elements", BenchmarkingPull(q, 100))
	q = ctl.NewHeapNodePriorityQueue()
	b.Run("Benchmark pull operation with preload of 100000 elements", BenchmarkingPull(q, 100000))
}

func BenchmarkingPush(q ctl.NodePriorityQueue, elements int) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i <= elements; i++ {
			q.Upsert(randomCapturer())
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			q.Upsert(randomCapturer())
		}
	}
}

func randomCapturer() *ctl.Node {
	index := random.Intn(len(nodeSet) - 1)
	uid := uuid.New().String()
	node := nodeSet[index]
	node.UUID = uid
	return node
}

func BenchmarkingPull(q ctl.NodePriorityQueue, elements int) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i <= elements; i++ {
			q.Upsert(randomCapturer())
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			q.Next()
		}
	}
}
