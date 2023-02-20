package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"order-service/internal/domain/dto"
)

type BrokerProducer struct {
	w kafka.Writer
}

type ProducerConfig struct {
	//Username string
	//Password string
	Brokers []string
}

func NewBrokerProducer(cfg ProducerConfig) *BrokerProducer {
	client := &BrokerProducer{}

	w := kafka.Writer{
		Addr: kafka.TCP(cfg.Brokers...),

		//Transport: &kafka.Transport{
		//	SASL: plain.Mechanism{
		//		Username: cfg.Username,
		//		Password: cfg.Password,
		//	},
		//},
	}

	client.w = w
	return client
}

func (client *BrokerProducer) SendMessage(ctx context.Context, topic string, command dto.CommandDTO) error {
	message := RequestCommandFromDTO(command)

	marshaledMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message: %v", err)
	}

	err = client.w.WriteMessages(ctx, kafka.Message{
		Value: marshaledMessage,
		Topic: topic,
	})
	if err != nil {
		return fmt.Errorf("could not send message: %v", err)
	}

	return nil
}
