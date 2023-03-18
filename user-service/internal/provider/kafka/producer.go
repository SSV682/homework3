package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	domain "user-service/internal/domain/models"
)

type MultipleErr []error

func (m *MultipleErr) Error() string {
	b := bytes.NewBufferString("")

	for i, err := range *m {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteString(err.Error())
	}

	return b.String()
}

type BrokerProducer struct {
	w kafka.Writer
}

type ProducerConfig struct {
	Brokers []string
}

func NewBrokerProducer(cfg ProducerConfig) *BrokerProducer {
	client := &BrokerProducer{}

	w := kafka.Writer{
		Addr: kafka.TCP(cfg.Brokers...),
	}

	client.w = w
	return client
}

func (client *BrokerProducer) SendCommand(ctx context.Context, user domain.User, topicsName []string) error {
	message := NewCommand(user)
	fmt.Sprintf("message: %#v", message)

	marshaledMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message: %v", err)
	}

	errs := MultipleErr{}

	for _, topic := range topicsName {

		err = client.w.WriteMessages(ctx, kafka.Message{
			Value: marshaledMessage,
			Headers: []kafka.Header{
				{
					Key:   "publisher",
					Value: []byte("billing_service"),
				},
			},
			Topic: topic,
		})
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errors.New(errs.Error())
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
