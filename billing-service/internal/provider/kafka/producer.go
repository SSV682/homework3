package kafka

import (
	domain "billing-service/internal/domain/models"
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
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

func (client *BrokerProducer) SendCommand(ctx context.Context, command domain.ReadyResponseCommand) error {
	err := client.w.WriteMessages(ctx, kafka.Message{
		Value: command.Command,
		Topic: command.Topic,
	})
	if err != nil {
		return fmt.Errorf("could not send message: %v", err)
	}

	return nil
}

//func (client *BrokerProducer) SendMessage(ctx context.Context, message domain.Message) error {
//	marshaledMessage, err := json.Marshal(message.Order)
//	if err != nil {
//		return fmt.Errorf("could not marshal message: %v", err)
//	}
//
//	err = client.w.WriteMessages(ctx, kafka.Message{
//		Value: marshaledMessage,
//		Topic: message.Topic,
//	})
//	if err != nil {
//		return fmt.Errorf("could not send message: %v", err)
//	}
//
//	return nil
//}
