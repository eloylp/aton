package ctl

import (
	"context"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/eloylp/aton/components/ctl/config"
	"github.com/eloylp/aton/components/ctl/metrics"
)

type Ctl struct {
	cfg            *config.Config
	metricsService *metrics.Service
	api            *http.Server
	logger         *logrus.Logger
	wg             *sync.WaitGroup
	L              *sync.Mutex
}

func (c *Ctl) Start() error {
	c.logger.Infof("starting CTL at %s", c.cfg.ListenAddress)
	c.initializeAPI()
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

func (c *Ctl) Shutdown() {
	c.logger.Info("started graceful shutdown sequence")
	// Close api server
	if err := c.api.Shutdown(context.TODO()); err != nil {
		c.logger.Errorf("ctl: shutdown: %v", err)
	}
	c.wg.Wait()
	c.logger.Infof("stopped CTL at %s", c.cfg.ListenAddress)
}
