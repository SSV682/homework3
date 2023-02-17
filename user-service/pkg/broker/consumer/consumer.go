package consumer

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"user-service/pkg/broker"
)

type Message struct {
	// Topic the broker topic for c message.
	Topic string

	// Value the actual message to use in broker.
	Value []byte
}

// Consumer representation a High-level broker Consumer instance.
type Consumer interface {
	// Consume joins a cluster of consumers for a given list of topics and
	// to get the messages&
	Consume() (*Message, error)

	// ConsumeContext is an analog of Consume but with context support.
	ConsumeContext(ctx context.Context) (*Message, error)

	// Commit fixes the offset to the last message read.
	Commit() error

	// CommitOffset commits the offset of partition
	CommitOffset(partition int32, offset int64) error

	// Close closes the producer, releasing any open resources.
	Close() error

	// Ping verifies a connection to the broker.
	Ping() error

	// CurrentOffset returns Offset from last message.
	CurrentOffset(partition int32) int64

	// LimitOffset returns high value for partition in topic.
	LimitOffset(partition int32) (int64, error)

	// Partitions returns the partition from which messages were received.
	Partitions() []int32
}

type consumer struct {
	kafkaConsumer       *kafka.Consumer
	pollMessageTimeout  time.Duration
	pingTimeout         time.Duration
	topicName           string
	partitionsOfOffsets map[int32]kafka.Offset
	mu                  sync.RWMutex
}

func NewConsumer(config *Config) (Consumer, error) {
	err := config.init()
	if err != nil {
		return nil, err
	}

	configMap := kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(config.BrokerAddresses, ", "),
		"group.id":           config.GroupId,
		"auto.offset.reset":  config.AutoOffsetReset.String(),
		"enable.auto.commit": !config.DisableAutoCommit,
		"security.protocol":  broker.SASLPlaintextSecurityProtocol,
		"sasl.mechanism":     broker.SASLMechanism,
		"sasl.username":      config.SASL.Username,
		"sasl.password":      config.SASL.Password,
	}

	if config.SSL != nil {
		configMap["security.protocol"] = broker.SASLSSLSecurityProtocol
		configMap["ssl.ca.location"] = config.SSL.CALocation
	}

	consumerClient, err := kafka.NewConsumer(&configMap)
	if err != nil {
		return nil, err
	}

	err = consumerClient.Subscribe(config.Topic, nil)
	if err != nil {
		return nil, err
	}

	return &consumer{
		kafkaConsumer:       consumerClient,
		pollMessageTimeout:  config.PollMessageTimeout,
		pingTimeout:         config.PingTimeout,
		topicName:           config.Topic,
		partitionsOfOffsets: make(map[int32]kafka.Offset),
	}, nil
}

func (c *consumer) Consume() (*Message, error) {
	for {
		msg, err := c.readMessage()
		if err != nil {
			return nil, err
		}

		if msg == nil {
			continue
		}

		c.setPartitionOffset(msg.TopicPartition.Partition, msg.TopicPartition.Offset)

		return &Message{
			Topic: *msg.TopicPartition.Topic,
			Value: msg.Value,
		}, nil
	}
}

func (c *consumer) ConsumeContext(ctx context.Context) (*Message, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			msg, err := c.readMessage()
			if err != nil {
				return nil, err
			}

			if msg == nil {
				continue
			}

			c.setPartitionOffset(msg.TopicPartition.Partition, msg.TopicPartition.Offset)

			return &Message{
				Topic: *msg.TopicPartition.Topic,
				Value: msg.Value,
			}, nil
		}
	}
}

func (c *consumer) Commit() error {
	_, err := c.kafkaConsumer.Commit()
	return err
}

func (c *consumer) CommitOffset(partition int32, offset int64) error {
	_, err := c.kafkaConsumer.CommitOffsets([]kafka.TopicPartition{
		{
			Topic:     &c.topicName,
			Partition: partition,
			Offset:    kafka.Offset(offset),
		},
	})

	return err
}

func (c *consumer) Close() error {
	return c.kafkaConsumer.Close()
}

func (c *consumer) Ping() error {
	_, err := c.kafkaConsumer.GetMetadata(nil, false, int(c.pingTimeout.Milliseconds()))
	return err
}

func (c *consumer) readMessage() (*kafka.Message, error) {
	ev := c.kafkaConsumer.Poll(int(c.pollMessageTimeout.Milliseconds()))
	if ev == nil {
		return nil, nil
	}

	switch e := ev.(type) {
	case *kafka.Message:
		if e.TopicPartition.Error != nil {
			return e, e.TopicPartition.Error
		}

		return e, nil
	case kafka.Error:
		return nil, e
	}

	// Ignore other event types
	return nil, nil
}

func (c *consumer) setPartitionOffset(partition int32, offset kafka.Offset) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.partitionsOfOffsets[partition] = offset
}

func (c *consumer) CurrentOffset(partition int32) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return int64(c.partitionsOfOffsets[partition])
}

func (c *consumer) LimitOffset(partition int32) (int64, error) {
	_, high, err := c.kafkaConsumer.GetWatermarkOffsets(c.topicName, partition)
	return high, err
}

func (c *consumer) Partitions() []int32 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	res := make([]int32, 0, len(c.partitionsOfOffsets))
	for key, _ := range c.partitionsOfOffsets {
		res = append(res, key)
	}

	return res
}
