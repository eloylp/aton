package ctl

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/eloylp/aton/internal/ctl/config"
	"github.com/eloylp/aton/internal/ctl/www"
	"github.com/eloylp/aton/internal/logging"
	"github.com/eloylp/aton/internal/proto"
)

type CapturerRegistry map[string]Capturer

type Ctl struct {
	cfg            *config.Config
	detectorClient DetectorClient
	detectorIn     chan<- *proto.RecognizeRequest
	detectorOut    <-chan *proto.RecognizeResponse
	capturers      CapturerRegistry
	stats          *Stats
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
		Handler:      www.Router(),
		ReadTimeout:  cfg.APIReadTimeout,
		WriteTimeout: cfg.APIWriteTimeout,
	}
	stats := &Stats{}
	stats.IncDetectors(int32(len(cfg.Detectors)))
	ctl := &Ctl{
		cfg:            cfg,
		detectorClient: dc,
		stats:          stats,
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
	var err error
	c.detectorIn, c.detectorOut, err = c.detectorClient.Recognize(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

func (c *Ctl) initializeResultProcessor() {
	c.wg.Add(1)
	go func() {
		for resp := range c.detectorOut {
			if resp.Success {
				c.stats.IncSuccess()
			} else {
				c.stats.IncFailed()
			}
			c.logger.Info("detected: " + strings.Join(resp.Names, ","))
		}
		c.wg.Done()
	}()
}

func (c *Ctl) Shutdown() error {
	// Close api server
	err := c.api.Shutdown(context.TODO())
	// Close capturers. Stop receiving more data to the system.
	for _, capt := range c.capturers {
		capt.Close()
	}
	// Close detectors client
	c.detectorClient.Shutdown()
	c.wg.Wait()
	return err
}

func (c *Ctl) Stats() *Stats {
	return c.stats
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
		go capturer.Start()
		for fr := range capturer.Output() {
			c.detectorIn <- &proto.RecognizeRequest{
				Image:     fr.Data,
				CreatedAt: timestamppb.New(fr.Timestamp),
			}
		}
		c.wg.Done()
	}(capt)
	c.stats.IncCapturers()
}
