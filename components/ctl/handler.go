package ctl

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type DetectorHandler struct {
	detector       *Detector
	client         DetectorClient
	priorityQueue  DetectorPriorityQueue
	processStopper chan struct{}
	wg             *sync.WaitGroup
	logger         *logrus.Logger
}

func NewDetectorHandler(detector *Detector, client DetectorClient, logger *logrus.Logger) *DetectorHandler {
	return &DetectorHandler{
		detector:       detector,
		client:         client,
		logger:         logger,
		wg:             &sync.WaitGroup{},
		processStopper: make(chan struct{}),
	}
}

func (ds *DetectorHandler) Start() error {
	if err := ds.client.Connect(); err != nil {
		return fmt.Errorf("ctl: could not connect to detector %s: %w", ds.detector.UUID, err)
	}
	ds.wg.Add(2)
	go ds.processStatus()
	go ds.processResults()
	return nil
}

func (ds *DetectorHandler) processResults() {
	defer ds.wg.Done()
	for {
		select {
		case <-ds.processStopper:
			ds.logger.Infof("detectorHandler: closed processing results of %s", ds.detector.UUID)
			return
		default:
			r, err := ds.client.NextResult()
			if err != nil {
				ds.logger.Errorf("detectorHandler: error obtaining next result: %s", err)
				return
			}
			resultFormat := "ctl: result: %s - %d (%s) - %d | %s | %s"
			ds.logger.Infof(resultFormat,
				r.DetectorUUID,
				len(r.Recognized),
				strings.Join(r.Recognized, ","),
				r.TotalEntities,
				r.CapturedAt,
				r.RecognizedAt)
		}
	}
}

func (ds *DetectorHandler) processStatus() {
	defer ds.wg.Done()
	for {
		select {
		case <-ds.processStopper:
			ds.logger.Infof("detectorHandler: closed processing status of %s", ds.detector.UUID)
			return
		default:
			s, err := ds.client.NextStatus()
			if err != nil {
				ds.logger.Errorf("detectorHandler: error obtaining next status: %s", err)
				return
			}
			ds.detector.Status = s
			ds.priorityQueue.Upsert(ds.detector)
		}
	}
}

func (ds *DetectorHandler) Shutdown() error {
	close(ds.processStopper)
	ds.wg.Wait()
	return nil
}
