package config

import (
	"fmt"
	"io"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	DefaultAPIReadTimeout  = 5 * time.Second
	DefaultAPIWriteTimeout = 5 * time.Second
)

type Config struct {
	ListenAddress   string   `split_words:"true" ,default:"0.0.0.0:8081"`
	Detector        Detector `split_words:"true" ,required:"true"`
	APIReadTimeout  time.Duration
	APIWriteTimeout time.Duration
	LogFormat       string `default:"json" ,split_words:"true"`
	LogOutput       io.Writer
}

type Option func(*Config)

func WithListenAddress(addr string) Option {
	return func(cfg *Config) {
		cfg.ListenAddress = addr
	}
}

func WithDetector(uuid, addr string) Option {
	return func(cfg *Config) {
		cfg.Detector = Detector{
			UUID:    uuid,
			Address: addr,
		}
	}
}

type Detector struct {
	UUID    string `split_words:"true" ,required:"true"`
	Address string `split_words:"true" ,required:"true"`
}

func FromEnv() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("ATON_CTL", cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return cfg, nil
}

func WithLogFormat(format string) Option {
	return func(cfg *Config) {
		cfg.LogFormat = format
	}
}

func WithLogOutput(w io.Writer) Option {
	return func(cfg *Config) {
		cfg.LogOutput = w
	}
}
