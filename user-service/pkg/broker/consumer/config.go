package consumer

import (
	"errors"
	"time"

	"user-service/pkg/broker"
)

// OffsetResetStrategy represents of the auto offset reset strategy
type OffsetResetStrategy int

const (
	Smallest OffsetResetStrategy = iota
	Earliest
	Beginning
	Largest
	Latest
	End
	Error
)

func (s *OffsetResetStrategy) IsValid() bool {
	switch *s {
	case Smallest,
		Earliest,
		Beginning,
		Largest,
		Latest,
		End,
		Error:
		return true
	default:
		return false
	}
}

func (s *OffsetResetStrategy) String() string {
	return [...]string{
		"smallest",
		"earliest",
		"beginning",
		"largest",
		"latest",
		"end",
		"error",
	}[*s]
}

// A Config specifies the configuration for the consumers.
type Config struct {
	*broker.NetworkConfig

	// GroupId client group id string.
	GroupId string

	// Topic the broker topic for s message.
	Topic string

	// AutoOffsetReset the action to take when there is no initial offset in offset store
	// or the desired offset is out of range.
	AutoOffsetReset OffsetResetStrategy

	// DisableAutoCommit when set to false, the consumer offsets
	// will periodically be fixed in the background.
	// If s property is set to true, the offsets are not fixed.
	DisableAutoCommit bool

	// PollMessageTimeout timeout during which it will be blocked for the next read.
	PollMessageTimeout time.Duration

	// PingTimeout a temporary measure for consumer.Ping() that is not blocked
	PingTimeout time.Duration
}

func (s *Config) init() error {
	if s.NetworkConfig == nil {
		return errors.New("NetworkConfig cannot be nil for broker consumer configuration")
	}

	if err := s.NetworkConfig.Init(); err != nil {
		return err
	}

	if len(s.GroupId) == 0 {
		return errors.New("group id cannot be empty for broker consumer configuration")
	}

	if len(s.Topic) == 0 {
		return errors.New("topic cannot be empty for broker consumer configuration")
	}

	if !s.AutoOffsetReset.IsValid() {
		return errors.New("auto offset reset strategy is not valid")
	}

	if s.PollMessageTimeout == 0 {
		s.PollMessageTimeout = 10 * time.Second
	}

	if s.PingTimeout == 0 {
		s.PingTimeout = 1 * time.Second
	}

	return nil
}
