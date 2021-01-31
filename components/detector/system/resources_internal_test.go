package system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetworkGauge(t *testing.T) {
	t1 := NetworkCounters{
		RxBytes: 10000,
		TxBytes: 10000,
	}
	t2 := NetworkCounters{
		RxBytes: 15000,
		TxBytes: 15000,
	}
	assert.Equal(t, uint64(5000), t2.minus(t1).RxBytesSec)
	assert.Equal(t, uint64(5000), t2.minus(t1).TxBytesSec)
}
