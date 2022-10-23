package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App        `yaml:"app"`
		HTTP       `yaml:"http"`
		Log        `yaml:"logger"`
		Connection `yaml:"conn"`
		Timeout    `yaml:"timeout"`
	}

	App struct {
		Name    string `env:"APP_NAME" env-required:"true" yaml:"name"`
		Version string `env:"APP_VERSION" env-required:"true" yaml:"version"`
	}

	HTTP struct {
		Port string `env:"APP_PORT" env-required:"true" yaml:"port"`
	}

	Log struct {
		Level string `env:"APP_LOGLEVEL" env-required:"true" yaml:"log_level"`
	}

	Connection struct {
		Dbname   string `env:"DB_NAME" env-required:"true" yaml:"dbname"`
		User     string `env:"DB_USER" env-required:"true" yaml:"user"`
		Password string `env:"DB_PASSWORD" env-required:"true" yaml:"password"`
		Host     string `env:"DB_HOST" env-required:"true" yaml:"host"`
		Port     string `env:"DB_PORT" env-required:"true" yaml:"port"`
	}

	Timeout struct {
		Duration string `env:"APP_TIMEOUT" env-required:"true" yaml:"duration"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
