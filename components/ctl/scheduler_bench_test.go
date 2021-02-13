package ctl_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/eloylp/aton/components/ctl"
)

var (
	random      = rand.New(rand.NewSource(time.Now().Unix()))
	detectorSet = [...]*ctl.Detector{
		LeastUtilizedDetector(),
		OneThirdUtilizedDetector(),
		MidUtilizedDetector(),
		AverageUtilizedDetector(),
		FullUtilizedDetector(),
	}
)

func Benchmark_HeapPriorityQueue(b *testing.B) {
	q := ctl.NewHeapDetectorPriorityQueue()
	b.Run("Benchmark push operation with preload of 10 elements", BenchmarkingPush(q, 10))
	q = ctl.NewHeapDetectorPriorityQueue()
	b.Run("Benchmark push operation with preload of 100 elements", BenchmarkingPush(q, 100))
	q = ctl.NewHeapDetectorPriorityQueue()
	b.Run("Benchmark push operation with preload of 100000 elements", BenchmarkingPush(q, 100000))

	q = ctl.NewHeapDetectorPriorityQueue()
	b.Run("Benchmark pull operation with preload of 10 elements", BenchmarkingPull(q, 10))
	q = ctl.NewHeapDetectorPriorityQueue()
	b.Run("Benchmark pull operation with preload of 100 elements", BenchmarkingPull(q, 100))
	q = ctl.NewHeapDetectorPriorityQueue()
	b.Run("Benchmark pull operation with preload of 100000 elements", BenchmarkingPull(q, 100000))

}

func BenchmarkingPush(q ctl.DetectorPriorityQueue, elements int) func(*testing.B) {
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

func randomCapturer() *ctl.Detector {
	index := random.Intn(len(detectorSet) - 1)
	uid := uuid.New().String()
	detector := detectorSet[index]
	detector.UUID = uid
	return detector
}

func BenchmarkingPull(q ctl.DetectorPriorityQueue, elements int) func(*testing.B) {
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
