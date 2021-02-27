package system

import (
	"fmt"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
)

func CPUCount() int {
	c, err := cpu.Get()
	if err != nil {
		panicErr(err)
	}
	return c.CPUCount
}

func LoadAverage() LoadAverageGauge {
	l, err := loadavg.Get()
	if err != nil {
		panicErr(err)
	}
	return LoadAverageGauge{
		LoadAvg1:  l.Loadavg1,
		LoadAvg5:  l.Loadavg5,
		LoadAvg15: l.Loadavg15,
	}
}

type LoadAverageGauge struct {
	LoadAvg1  float64
	LoadAvg5  float64
	LoadAvg15 float64
}

func Memory() MemoryGauge {
	m, err := memory.Get()
	if err != nil {
		panicErr(err)
	}
	return MemoryGauge{
		TotalBytes: m.Total,
		UsedBytes:  m.Total - m.Used,
	}
}

type MemoryGauge struct {
	TotalBytes, UsedBytes uint64
}

func panicErr(err error) {
	panic(fmt.Errorf("system: resources: %w", err))
}

func Network() NetworkGauge {
	n1 := NetworkCount()
	time.Sleep(time.Second)
	n2 := NetworkCount()
	return n2.minus(n1)
}

// NetworkCount tries to calculate the most used interface and
// return only its usage count.
func NetworkCount() NetworkCounters {
	n, err := network.Get()
	if err != nil {
		panicErr(err)
	}
	var rxBytes uint64
	var txBytes uint64
	for _, nn := range n {
		if rxBytes < nn.RxBytes {
			rxBytes = nn.RxBytes
			txBytes = nn.TxBytes
		}
	}
	return NetworkCounters{
		RxBytes: rxBytes,
		TxBytes: txBytes,
	}
}

type NetworkGauge struct {
	RxBytesSec, TxBytesSec uint64
}

type NetworkCounters struct {
	RxBytes, TxBytes uint64
}

func (n NetworkCounters) minus(previousMeasure NetworkCounters) NetworkGauge {
	return NetworkGauge{
		RxBytesSec: n.RxBytes - previousMeasure.RxBytes,
		TxBytesSec: n.TxBytes - previousMeasure.TxBytes,
	}
}
