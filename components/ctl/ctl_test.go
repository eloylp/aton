package ctl_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/eloylp/aton/components/ctl"
	"github.com/eloylp/aton/components/ctl/metrics"
)

func TestCtl_AddCapturer_WithoutDetector(t *testing.T) {
	ms := metrics.NewService()
	logOutput := bytes.NewBuffer(nil)
	logger := logrus.New()
	logger.SetOutput(logOutput)
	c := ctl.NewCtl(logger, ms, func(addr string, l *logrus.Logger, service *metrics.Service) ctl.DetectorClient {
		return NewFakeDetectorClient(addr, l, service)
	})
	ctx := context.Background()
	err := c.AddCapturer(ctx, "UUID", "http://example.com")
	assert.EqualError(t, err, "ctl: cannot find suitable detector")

	logString := logOutput.String()
	assert.Contains(t, logString, `level=error`)
	assert.Contains(t, logString, `ctl: not suitable node for capturer UUID with URL http://example.com`)
}

func TestCtl_ExecutionAndShutdown(t *testing.T) {
	ms := metrics.NewService()
	logOutput := bytes.NewBuffer(nil)
	logger := logrus.New()
	logger.SetOutput(logOutput)

	dc := NewFakeDetectorClient("0.0.0.0:8080", logger, ms)
	dc.On("Connect").Return(nil)
	dc.On("NextStatus").Return(LeastUtilizedDetector().Status, nil)
	fixeNow, _ := time.Parse("2006", "2021")
	dc.On("NextResult").Return(&ctl.Result{
		DetectorUUID:  "UUID",
		Recognized:    []string{"alice", "bob"},
		TotalEntities: 15,
		RecognizedAt:  fixeNow.Add(time.Second),
		CapturedAt:    fixeNow,
	}, nil)
	c := ctl.NewCtl(logger, ms, func(addr string, l *logrus.Logger, service *metrics.Service) ctl.DetectorClient {
		return dc
	})
	_, err := c.AddDetector("0.0.0.0:8080")
	assert.NoError(t, err)
	time.Sleep(30 * time.Millisecond) // wait for goroutines.
	err = c.Shutdown(context.Background())
	assert.NoError(t, err)
	dc.AssertExpectations(t)

	logString := logOutput.String()
	assert.NotContains(t, logString, `level=error`)
	assert.Contains(t, logString, `ctl: added detector at 0.0.0.0:8080 for`)
	assert.Contains(t, logString, `ctl: result: UUID - 2 (alice,bob) - 15 | 2021-01-01 00:00:00 +0000 UTC | 2021-01-01 00:00:01 +0000 UTC`)
	assert.Contains(t, logString, `detectorHandler: closed processing status of`)
	assert.Contains(t, logString, `detectorHandler: closed processing results of`)
}

type FakeDetectorClient struct {
	addr    string
	l       *logrus.Logger
	service *metrics.Service
	mock.Mock
}

func NewFakeDetectorClient(addr string, l *logrus.Logger, service *metrics.Service) *FakeDetectorClient {
	return &FakeDetectorClient{
		addr:    addr,
		l:       l,
		service: service,
	}
}

func (d *FakeDetectorClient) Connect() error {
	return d.Called().Error(0)
}

func (d *FakeDetectorClient) LoadCategories(ctx context.Context, r *ctl.LoadCategoriesRequest) error {
	panic("implement me")
}

func (d *FakeDetectorClient) AddCapturer(ctx context.Context, r *ctl.AddCapturerRequest) error {
	args := d.Called(ctx, r)
	return args.Error(0)
}

func (d *FakeDetectorClient) RemoveCapturer(ctx context.Context, r *ctl.RemoveCapturerRequest) error {
	panic("implement me")
}

func (d *FakeDetectorClient) NextStatus() (*ctl.Status, error) {
	args := d.Called()
	return args.Get(0).(*ctl.Status), args.Error(1)
}

func (d *FakeDetectorClient) NextResult() (*ctl.Result, error) {
	args := d.Called()
	return args.Get(0).(*ctl.Result), args.Error(1)
}

func (d *FakeDetectorClient) Shutdown() error {
	panic("implement me")
}
