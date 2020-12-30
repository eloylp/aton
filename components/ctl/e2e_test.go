// +build e2e

package ctl_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/eloylp/aton/pkg/test/helper"
	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/components/ctl"
	"github.com/eloylp/aton/components/ctl/config"
)

func TestStartStopSequence(t *testing.T) {
	logOutput := bytes.NewBuffer(nil)
	c, err := ctl.New(
		config.WithListenAddress("0.0.0.0:10001"),
		config.WithDetector("test-detector", "127.0.0.1:10002"),
		config.WithLogOutput(logOutput),
		config.WithLogFormat("text"),
	)
	assert.NoError(t, err)
	go c.Start()
	helper.TryConnectTo(t, "0.0.0.0:10001", time.Second)
	c.Shutdown()

	logO := logOutput.String()
	assert.Contains(t, logO, "starting CTL at 0.0.0.0:10001")
	assert.Contains(t, logO, "started graceful shutdown sequence")
	assert.Contains(t, logO, "stopped CTL at 0.0.0.0:10001")
}
