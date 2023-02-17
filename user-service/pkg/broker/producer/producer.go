package producer

import (
	"strings"
	"time"
	"user-service/pkg/broker"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Producer publishes broker messages, blocking until they have been acknowledged.
// It routes messages to the correct broker, refreshing metadata as appropriate,
// and parses responses for errors.
type Producer interface {
	// Produce produces a given message, and returns only when it either has
	// succeeded or failed to produce. It will return an error if the message failed to produce.
	Produce(message *broker.Message) error

	// Close closes the producer, releasing any open resources.
	Close() error

	// Ping verifies a connection to the broker.
	Ping() error
}

type producer struct {
	kafkaProducer *kafka.Producer
	pingTimeout   time.Duration
}

func NewProducer(config *Config) (Producer, error) {
	err := config.init()
	if err != nil {
		return nil, err
	}

	configMap := kafka.ConfigMap{
		"bootstrap.servers":   strings.Join(config.BrokerAddresses, ", "),
		"security.protocol":   broker.SASLPlaintextSecurityProtocol,
		"sasl.mechanism":      broker.SASLMechanism,
		"sasl.username":       config.SASL.Username,
		"sasl.password":       config.SASL.Password,
		"go.delivery.reports": false,
	}

	if config.SSL != nil {
		configMap["security.protocol"] = broker.SASLSSLSecurityProtocol
		configMap["ssl.ca.location"] = config.SSL.CALocation
	}

	kafkaProducer, err := kafka.NewProducer(&configMap)
	if err != nil {
		return nil, err
	}

	return &producer{
		kafkaProducer: kafkaProducer,
		pingTimeout:   config.PingTimeout,
	}, nil
}

func (p *producer) Produce(message *broker.Message) error {
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err := p.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &message.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: message.Value,
	}, deliveryChan)
	if err != nil {
		return err
	}

	event := <-deliveryChan
	kafkaMessage := event.(*kafka.Message)
	if kafkaMessage.TopicPartition.Error != nil {
		return kafkaMessage.TopicPartition.Error
	}

	return nil
}

func (p *producer) Close() error {
	p.kafkaProducer.Close()

	return nil
}

func (p *producer) Ping() error {
	_, err := p.kafkaProducer.GetMetadata(nil, false, int(p.pingTimeout.Milliseconds()))
	return err
}
