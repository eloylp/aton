package detector

import (
	"fmt"
	"io"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/eloylp/aton/components/detector/metrics"
	"github.com/eloylp/aton/components/video"
)

type CapturerStatus struct {
	UUID   string
	URL    string
	Status string
}
type CapturerHandler struct {
	capturers      map[string]Capturer
	logger         *logrus.Logger
	output         chan *video.Capture
	metricsService *metrics.Service
	wg             sync.WaitGroup
	L              sync.RWMutex
}

func NewCapturerHandler(
	logger *logrus.Logger,
	metricsService *metrics.Service,
	backboneBuffSize int,
) *CapturerHandler {
	return &CapturerHandler{
		logger:         logger,
		capturers:      map[string]Capturer{},
		metricsService: metricsService,
		output:         make(chan *video.Capture, backboneBuffSize),
	}
}

func (th *CapturerHandler) AddCapturer(t Capturer) {
	th.L.Lock()
	defer th.L.Unlock()
	th.capturers[t.UUID()] = t
	th.logger.Infof("capturerHandler: added target with UUID: %s", t.UUID())
	th.initializeCapturer(t)
}

func (th *CapturerHandler) AddMJPEGCapturer(uuid, url string, maxFrameBuffer int) error {
	capt, err := video.NewMJPEGCapturer(uuid, url, maxFrameBuffer, th.logger)
	if err != nil {
		return err
	}
	th.AddCapturer(capt)
	return nil
}

func (th *CapturerHandler) initializeCapturer(t Capturer) {
	th.logger.Infof("capturerHandler: starting target with UUID: %s", t.UUID())
	th.wg.Add(1)
	go func() {
		th.metricsService.CapturerUP(t.UUID())
		defer th.metricsService.CapturerDown(t.UUID())
		go t.Start()
		for {
			fr, err := t.NextOutput()
			if err == io.EOF {
				break
			}
			th.metricsService.IncCapturerReceivedFramesTotal(t.UUID())
			if err != nil {
				th.metricsService.IncCapturerFailedFramesTotal(t.UUID())
				th.logger.Errorf("capturerHandler: target: %v", err)
				continue
			}
			th.output <- fr
		}
		th.wg.Done()
	}()
}

func (th *CapturerHandler) BackboneLen() int {
	return len(th.output)
}

func (th *CapturerHandler) Shutdown() {
	th.L.Lock()
	defer th.L.Unlock()
	for _, t := range th.capturers {
		th.logger.Infof("capturerHandler: closing target with UUID: %s", t.UUID())
		t.Close()
	}
	th.wg.Wait()
}

func (th *CapturerHandler) Status() []CapturerStatus {
	th.L.RLock()
	defer th.L.RUnlock()
	var status []CapturerStatus
	for _, t := range th.capturers {
		status = append(status, CapturerStatus{
			UUID:   t.UUID(),
			Status: t.Status(),
			URL:    t.TargetURL(),
		})
	}
	return status
}

func (th *CapturerHandler) RemoveCapturer(uuid string) error {
	th.L.Lock()
	defer th.L.Unlock()
	if _, ok := th.capturers[uuid]; !ok {
		return fmt.Errorf("capturerHandler: capturer with UUID %s not found", uuid)
	}
	th.capturers[uuid].Close()
	delete(th.capturers, uuid)
	return nil
}
