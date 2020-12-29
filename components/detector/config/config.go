package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Option func(cfg *Config)

type Config struct {
	UUID              string `split_words:"true" ,required:"true"`
	ListenAddr        string `split_words:"true" ,default:"0.0.0.0:8080"`
	MetricsListenAddr string `split_words:"true" ,default:"0.0.0.0:8081"`
	ModelDir          string `split_words:"true" ,required:"true"`
	LogFormat         string `default:"json" ,split_words:"true"`
}

func WithUUID(uuid string) Option {
	return func(cfg *Config) {
		cfg.UUID = uuid
	}
}

func WithListenAddress(addr string) Option {
	return func(cfg *Config) {
		cfg.ListenAddr = addr
	}
}

func WithMetricsAddress(addr string) Option {
	return func(cfg *Config) {
		cfg.MetricsListenAddr = addr
	}
}

func WithModelDir(dir string) Option {
	return func(cfg *Config) {
		cfg.ModelDir = dir
	}
}

func WithLogFormat(logFormat string) Option {
	return func(cfg *Config) {
		cfg.LogFormat = logFormat
	}
}

func FromEnv() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("DETECTOR", cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return cfg, nil
}
