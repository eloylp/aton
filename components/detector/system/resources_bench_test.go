package system_test

import (
	"testing"

	"github.com/eloylp/aton/components/detector/system"
)

func BenchmarkCPUCount(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		system.CPUCount()
	}
}

func BenchmarkLoadAverage(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		system.LoadAverage()
	}
}

func BenchmarkMemory(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		system.Memory()
	}
}

func BenchmarkNetwork(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		system.NetworkCount()
	}
}
