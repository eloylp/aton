// +build e2e

package detector_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/components/detector"
	"github.com/eloylp/aton/components/detector/config"
)

func TestStartStopSequence(t *testing.T) {
	logOutput := bytes.NewBuffer(nil)
	d, err := detector.New(
		config.WithListenAddress("0.0.0.0:10002"),
		config.WithMetricsAddress("0.0.0.0:10003"),
		config.WithLogOutput(logOutput),
		config.WithLogFormat("text"),
		config.WithModelDir("../../models"),
	)
	assert.NoError(t, err)

	go d.Start()
	time.Sleep(time.Second) // Wait for initialization
	d.Shutdown()

	logO := logOutput.String()
	assert.Contains(t, logO, "starting detector service at 0.0.0.0:10002")
	assert.Contains(t, logO, "starting detector metrics at 0.0.0.0:10003")
	assert.Contains(t, logO, "gracefully shutdown started.")
	assert.Contains(t, logO, "stopped detector service at 0.0.0.0:10002")
	assert.Contains(t, logO, "stopped detector metrics at 0.0.0.0:10003")
	assert.NotContains(t, logO, "level=error")

	t.Log(logO)
}
