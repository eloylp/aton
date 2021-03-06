// +build integration

package system_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/components/node/system"
)

func TestCPUCount(t *testing.T) {
	assert.Greater(t, system.CPUCount(), 1)
}

func TestLoadAverage(t *testing.T) {
	loadAverage := system.LoadAverage()
	assert.Greater(t, loadAverage.LoadAvg1, float64(0))
	assert.Greater(t, loadAverage.LoadAvg5, float64(0))
	assert.Greater(t, loadAverage.LoadAvg15, float64(0))
}

func TestMem(t *testing.T) {
	mem := system.Memory()
	assert.Greater(t, mem.UsedBytes, uint64(0))
	assert.Greater(t, mem.TotalBytes, uint64(0))
	assert.Greater(t, mem.TotalBytes, mem.UsedBytes)
}
