package ctl

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/eloylp/aton/internal/ctl/config"
	"github.com/eloylp/aton/internal/ctl/metrics"
	"github.com/eloylp/aton/internal/ctl/www"
	"github.com/eloylp/aton/internal/proto"
)

type CapturerRegistry map[string]Capturer

type Ctl struct {
	cfg            *config.Config
	detectorClient DetectorClient
	capturers      CapturerRegistry
	metricsService *metrics.Service
	api            *http.Server
	logger         *logrus.Logger
	wg             *sync.WaitGroup
	L              *sync.Mutex
}

func New(dc DetectorClient, metricsService *metrics.Service, logger *logrus.Logger, opts ...config.Option) *Ctl {
	cfg := &config.Config{
		APIReadTimeout:  config.DefaultAPIReadTimeout,
		APIWriteTimeout: config.DefaultAPIWriteTimeout,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	api := &http.Server{
		Addr:         cfg.ListenAddress,
		Handler:      www.Router(metricsService.HTTPHandler()),
		ReadTimeout:  cfg.APIReadTimeout,
		WriteTimeout: cfg.APIWriteTimeout,
	}
	for _, d := range cfg.Detectors {
		metricsService.DetectorUP(d.UUID)
	}
	ctl := &Ctl{
		cfg:            cfg,
		detectorClient: dc,
		metricsService: metricsService,
		capturers:      CapturerRegistry{},
		api:            api,
		logger:         logger,
		wg:             &sync.WaitGroup{},
		L:              &sync.Mutex{},
	}
	return ctl
}

func (c *Ctl) Start() error {
	c.initializeAPI()
	if err := c.initializeDetectorClient(); err != nil {
		c.logger.Errorf("ctl: %v", err)
		return err
	}
	c.initializeResultProcessor()
	c.wg.Wait()
	return nil
}

func (c *Ctl) initializeAPI() {
	c.wg.Add(1)
	go func() {
		if err := c.api.ListenAndServe(); err != http.ErrServerClosed {
			c.logger.Errorf("ctl: %v", err)
		}
		c.wg.Done()
	}()
}

func (c *Ctl) initializeDetectorClient() error {
	if err := c.detectorClient.Connect(); err != nil {
		return err
	}
	if err := c.detectorClient.StartRecognize(context.TODO()); err != nil {
		return err
	}
	return nil
}

func (c *Ctl) initializeResultProcessor() {
	c.wg.Add(1)
	go func() {
		for {
			resp, err := c.detectorClient.NextRecognizeResponse()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.logger.Errorf("ctl: processor: %v", err)
				continue
			}
			if resp.Success {
				c.metricsService.IncProcessedFramesTotal(resp.ProcessedBy)
				if len(resp.Names) > 0 {
					c.logger.Info("initializeResultProcessor(): detected: " + strings.Join(resp.Names, ","))
				} else {
					c.metricsService.IncUnrecognizedFramesTotal(resp.ProcessedBy)
					c.logger.Info("initializeResultProcessor(): not detected: " + resp.Message)
				}
			} else {
				c.metricsService.IncFailedFramesTotal(resp.ProcessedBy)
			}
		}
		c.wg.Done()
	}()
}

func (c *Ctl) Shutdown() {
	// Close api server
	if err := c.api.Shutdown(context.TODO()); err != nil {
		c.logger.Errorf("ctl: shutdown: %v", err)
	}
	// Close capturers. Stop receiving more data to the system.
	for _, capt := range c.capturers {
		capt.Close()
	}
	// Close detectors client
	if err := c.detectorClient.Shutdown(); err != nil {
		c.logger.Errorf("ctl: shutdown: %v", err)
	}
	c.wg.Wait()
}

func (c *Ctl) AddCapturer(capt Capturer) {
	c.L.Lock()
	defer c.L.Unlock()
	c.capturers[capt.UUID()] = capt
	c.initializeCapturer(capt)
}

func (c *Ctl) initializeCapturer(capt Capturer) {
	c.wg.Add(1)
	go func(capturer Capturer) {
		c.metricsService.CapturerUP(capt.UUID())
		defer c.metricsService.CapturerDown(capt.UUID())
		go capturer.Start()
		for {
			fr, err := capturer.NextOutput()
			if err == io.EOF {
				break
			}
			c.metricsService.IncCapturerReceivedFramesTotal(capt.UUID())
			if err != nil {
				c.metricsService.IncCapturerFailedFramesTotal(capt.UUID())
				c.logger.Error("ctl: capturer: %w", err)
				continue
			}
			if err = c.detectorClient.SendToRecognize(&proto.RecognizeRequest{
				Image:     fr.Data,
				CreatedAt: timestamppb.New(fr.Timestamp),
			}); err != nil {
				c.logger.Error("ctl: capturer: sending: %w", err)
			}
		}
		c.wg.Done()
	}(capt)
}
