package kafka

import (
	"encoding/json"
	"fmt"
	"order-service/internal/domain/dto"
	"order-service/pkg/broker"
	"order-service/pkg/broker/producer"
)

type kafkaProducerProvider struct {
	kafkaProducer producer.Producer
}

func NewProvider(kafka *producer.Producer) (*kafkaProducerProvider, error) {
	producerService := &kafkaProducerProvider{
		kafkaProducer: *kafka,
	}

	return producerService, nil
}

func (p *kafkaProducerProvider) SendMessage(topic string, message dto.CommandDTO) error {
	marshaledMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	bm := broker.Message{
		Topic: topic,
		Value: marshaledMessage,
	}

	if err = p.kafkaProducer.Produce(&bm); err != nil {
		return fmt.Errorf("failed send message: %v", err)
	}
	return nil
}
