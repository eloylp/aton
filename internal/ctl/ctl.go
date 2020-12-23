package ctl

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/eloylp/aton/internal/ctl/config"
	"github.com/eloylp/aton/internal/ctl/metrics"
	"github.com/eloylp/aton/internal/ctl/www"
	"github.com/eloylp/aton/internal/logging"
	"github.com/eloylp/aton/internal/proto"
)

type CapturerRegistry map[string]Capturer

type Ctl struct {
	cfg            *config.Config
	detectorClient DetectorClient
	capturers      CapturerRegistry
	register       *metrics.Register
	api            *http.Server
	logger         logging.Logger
	wg             *sync.WaitGroup
	L              *sync.Mutex
}

func New(dc DetectorClient, opts ...config.Option) *Ctl {
	cfg := &config.Config{
		APIReadTimeout:  config.DefaultAPIReadTimeout,
		APIWriteTimeout: config.DefaultAPIWriteTimeout,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	api := &http.Server{
		Addr:         cfg.ListenAddress,
		Handler:      www.Router(metrics.NewHTTPHandler()),
		ReadTimeout:  cfg.APIReadTimeout,
		WriteTimeout: cfg.APIWriteTimeout,
	}
	register := metrics.NewRegister()
	for _, d := range cfg.Detectors {
		register.DetectorUP(d.UUID)
	}
	ctl := &Ctl{
		cfg:            cfg,
		detectorClient: dc,
		register:       register,
		capturers:      CapturerRegistry{},
		api:            api,
		logger:         logging.NewBasicLogger(cfg.LoggerOutput),
		wg:             &sync.WaitGroup{},
		L:              &sync.Mutex{},
	}
	return ctl
}

func (c *Ctl) Start() error {
	c.initializeAPI()
	if err := c.initializeDetectorClient(); err != nil {
		c.logger.Errorf("ctl: %w", err)
		return err
	}
	c.initializeResultProcessor()
	return nil
}

func (c *Ctl) initializeAPI() {
	c.wg.Add(1)
	go func() {
		if err := c.api.ListenAndServe(); err != http.ErrServerClosed {
			c.logger.Errorf("ctl: %w", err)
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
				c.logger.Errorf("ctl: processor: %w", err)
				continue
			}
			c.register.IncProcessedFramesTotal(resp.ProcessedBy)
			if resp.Success {
				c.logger.Info("detected: " + strings.Join(resp.Names, ","))
			} else {
				c.register.IncFailedFramesTotal(resp.ProcessedBy)
				c.logger.Info("not detected: " + resp.Message)
			}
		}
		c.wg.Done()
	}()
}

func (c *Ctl) Shutdown() {
	// Close api server
	if err := c.api.Shutdown(context.TODO()); err != nil {
		c.logger.Errorf("ctl: shutdown: %w", err)
	}
	// Close capturers. Stop receiving more data to the system.
	for _, capt := range c.capturers {
		capt.Close()
	}
	// Close detectors client
	if err := c.detectorClient.Shutdown(); err != nil {
		c.logger.Errorf("ctl: shutdown: %w", err)
	}
	c.wg.Wait()
}

func (c *Ctl) AddCapturer(capt Capturer) error {
	c.L.Lock()
	defer c.L.Unlock()
	c.capturers[capt.UUID()] = capt
	c.initializeCapturer(capt)
	return nil
}

func (c *Ctl) initializeCapturer(capt Capturer) {
	c.wg.Add(1)
	go func(capturer Capturer) {
		c.register.CapturerUP(capt.UUID())
		defer c.register.CapturerDown(capt.UUID())
		go capturer.Start()
		for {
			fr, err := capturer.NextOutput()
			if err == io.EOF {
				break
			}
			c.register.IncCapturerReceivedFramesTotal(capt.UUID())
			if err != nil {
				c.register.IncCapturerFailedFramesTotal(capt.UUID())
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