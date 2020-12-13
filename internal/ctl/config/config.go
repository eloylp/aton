package config

import (
	"io"
	"time"
)

const (
	DefaultAPIReadTimeout  = 5 * time.Second
	DefaultAPIWriteTimeout = 5 * time.Second
)

type Config struct {
	ListenAddress   string
	LoggerOutput    io.Writer
	Detectors       []Detector
	APIReadTimeout  time.Duration
	APIWriteTimeout time.Duration
}

type Option func(*Config)

func WithListenAddress(addr string) Option {
	return func(cfg *Config) {
		cfg.ListenAddress = addr
	}
}

func WithLoggerOutput(w io.Writer) Option {
	return func(cfg *Config) {
		cfg.LoggerOutput = w
	}
}

func WithDetectors(d ...Detector) Option {
	return func(cfg *Config) {
		cfg.Detectors = d
	}
}

type Detector struct {
	Address string
	UUID    string
}
