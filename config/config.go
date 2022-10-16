package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"`
		Version string `env-required:"true" yaml:"version"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
