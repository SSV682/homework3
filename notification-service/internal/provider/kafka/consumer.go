package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	domain "notification-service/internal/domain/models"
)

type BrokerConsumer struct {
	reader *kafka.Reader
}

func NewBrokerConsumer(brokers []string, topic, groupID string) *BrokerConsumer {
	client := &BrokerConsumer{}
	client.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})
	return client
}

func (c *BrokerConsumer) read(ctx context.Context) ([]byte, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return msg.Value, nil
}

func (c *BrokerConsumer) StartConsume(ctx context.Context, errCh chan error) (<-chan domain.Order, error) {
	payloadCh := make(chan domain.Order)

	go func() {
		defer close(payloadCh)

		for {
			select {
			case <-ctx.Done():
				log.Debug("Got context done! Closing consumer...")

				if err := c.reader.Close(); err != nil {
					errCh <- fmt.Errorf("consumer close: %v", err)
				}
				return

			default:
				message, err := c.consumeContext(ctx)
				if err != nil {
					errCh <- fmt.Errorf("consume message: %v", err)
					continue
				}

				var command Order

				if err = json.Unmarshal(message, &command); err != nil {
					errCh <- fmt.Errorf("unmarshal message: %v", err)
					continue
				}

				c := command.ToModel()
				payloadCh <- c
			}
		}
	}()

	return payloadCh, nil
}

func (c *BrokerConsumer) StartConsumeUserUpdate(ctx context.Context, errCh chan error) (<-chan domain.User, error) {
	payloadCh := make(chan domain.User)

	go func() {
		defer close(payloadCh)

		for {
			select {
			case <-ctx.Done():
				log.Debug("Got context done! Closing consumer...")

				if err := c.reader.Close(); err != nil {
					errCh <- fmt.Errorf("consumer close: %v", err)
				}
				return

			default:
				message, err := c.consumeContext(ctx)
				if err != nil {
					errCh <- fmt.Errorf("consume message: %v", err)
					continue
				}

				var command User

				if err = json.Unmarshal(message, &command); err != nil {
					errCh <- fmt.Errorf("unmarshal message: %v", err)
					continue
				}

				c := command.ToModel()
				payloadCh <- c
			}
		}
	}()

	return payloadCh, nil
}

func (c *BrokerConsumer) consumeContext(ctx context.Context) ([]byte, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			msg, err := c.read(ctx)
			if err != nil {
				return nil, err
			}

			if msg == nil {
				continue
			}

			return msg, nil
		}
	}
}
