package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/pkg/broker"
	"order-service/pkg/broker/consumer"

	log "github.com/sirupsen/logrus"
	"order-service/internal/domain/dto"
)

const groupId = "eventGroupStore"

type kafkaConsumerProvider struct {
	topic    string
	cfg      *broker.NetworkConfig
	consumer consumer.Consumer
}

func NewKafkaConsumerProvider(cfg *broker.NetworkConfig, topic string) *kafkaConsumerProvider {
	return &kafkaConsumerProvider{
		cfg:   cfg,
		topic: topic,
	}
}

func (e *kafkaConsumerProvider) StartConsume(ctx context.Context) (<-chan dto.OrderCommandDTO, <-chan error, error) {
	client, err := consumer.NewConsumer(&consumer.Config{
		NetworkConfig:   e.cfg,
		GroupId:         groupId,
		Topic:           e.topic,
		AutoOffsetReset: consumer.Latest,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("consumer create: %v", err)
	}

	e.consumer = client
	payloadCh, errCh := e.consume(ctx)

	return payloadCh, errCh, nil
}

func (e *kafkaConsumerProvider) consume(ctx context.Context) (<-chan dto.OrderCommandDTO, <-chan error) {
	payloadCh := make(chan dto.OrderCommandDTO)
	errCh := make(chan error)

	go func() {
		defer close(payloadCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				log.Debug("Got context done! Closing consumer...")

				if err := e.consumer.Close(); err != nil {
					errCh <- fmt.Errorf("consumer close: %v", err)
				}

				return
			default:
				message, err := e.consumer.ConsumeContext(ctx)
				if err != nil {
					errCh <- fmt.Errorf("consume message: %v", err)
					continue
				}

				var command dto.OrderCommandDTO

				if err = json.Unmarshal(message.Value, &command); err != nil {
					errCh <- fmt.Errorf("unmarshal message: %v", err)
					continue
				}

				log.Debugf("Log command message came: %v", command)

				payloadCh <- command
			}
		}
	}()

	return payloadCh, errCh
}
