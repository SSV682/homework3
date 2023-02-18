package producer

import (
	"errors"
	"order-service/pkg/broker"
	"time"
)

// A Config specifies the configuration for the producers.
type Config struct {
	*broker.NetworkConfig

	// PingTimeout a temporary measure for consumer.Ping() that is not blocked
	PingTimeout time.Duration
}

func (c *Config) init() error {
	if c.NetworkConfig == nil {
		return errors.New("NetworkConfig cannot be nil for broker producer configuration")
	}

	if err := c.NetworkConfig.Init(); err != nil {
		return err
	}

	if c.PingTimeout == 0 {
		c.PingTimeout = 1 * time.Second
	}

	return nil
}
