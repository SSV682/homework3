package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"order-service/internal/domain/dto"
)

const groupId = "eventGroupStore"

type BrokerConsumer struct {
	reader *kafka.Reader
}

func NewBrokerConsumer(brokers []string, topic string) *BrokerConsumer {
	client := &BrokerConsumer{}
	client.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupId,
		Topic:   topic,
	})
	return client
}

func (c *BrokerConsumer) Read(ctx context.Context) ([]byte, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return msg.Value, nil
}

func (c *BrokerConsumer) StartConsume(ctx context.Context) (<-chan dto.OrderCommandDTO, <-chan error, error) {
	payloadCh := make(chan dto.OrderCommandDTO)
	errCh := make(chan error)

	go func() {
		defer close(payloadCh)
		defer close(errCh)

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

				var command dto.OrderCommandDTO

				if err = json.Unmarshal(message, &command); err != nil {
					errCh <- fmt.Errorf("unmarshal message: %v", err)
					continue
				}

				payloadCh <- command
			}
		}
	}()

	return payloadCh, errCh, nil
}

func (c *BrokerConsumer) consumeContext(ctx context.Context) ([]byte, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			msg, err := c.Read(ctx)
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
