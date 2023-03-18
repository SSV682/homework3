package kafka

import (
	domain "billing-service/internal/domain/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"time"
)

const groupID = "eventGroupStore"

type BrokerConsumer struct {
	reader *kafka.Reader
}

func NewBrokerConsumer(brokers []string, topic string) *BrokerConsumer {
	client := &BrokerConsumer{}
	client.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:           brokers,
		GroupID:           groupID,
		Topic:             topic,
		HeartbeatInterval: 10 * time.Second,
	})
	return client
}

func (c *BrokerConsumer) read(ctx context.Context) (*kafka.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (c *BrokerConsumer) StartConsume(ctx context.Context) (<-chan domain.RequestCommand, <-chan domain.Account, <-chan error, error) {
	payloadCh := make(chan domain.RequestCommand)
	payloadUserCh := make(chan domain.Account)

	errCh := make(chan error)

	go func() {
		defer close(payloadCh)
		defer close(payloadUserCh)
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
					fmt.Println(err)
					errCh <- fmt.Errorf("consume message: %v", err)
					continue
				}

				fmt.Printf("message: %#v", message)
				var messageType string
				for _, v := range message.Headers {
					fmt.Printf("key: %s", v.Key)
					if v.Key == "publisher" {
						messageType = string(v.Value)
					}
				}

				switch messageType {
				case "billing_service":
					command, err := getBillingCommand(message.Value)
					if err != nil {
						continue
					}
					fmt.Printf("message come")
					payloadCh <- *command
				case "order_service":
					getUserCommand(message.Value)
					command, err := getUserCommand(message.Value)
					if err != nil {
						continue
					}
					payloadUserCh <- *command
				default:

				}

				return
			}
		}
	}()

	return payloadCh, payloadUserCh, errCh, nil
}

func getBillingCommand(message []byte) (*domain.RequestCommand, error) {
	var command BillingRequestCommand

	if err := json.Unmarshal(message, &command); err != nil {
		return nil, fmt.Errorf("unmarshal message: %v", err)
	}

	c, err := command.ToModel()
	if err != nil {
		return nil, fmt.Errorf("unmarshal message: %v", err)
	}

	return &c, nil
}

func getUserCommand(message []byte) (*domain.Account, error) {
	var command UserRequestCommand

	if err := json.Unmarshal(message, &command); err != nil {
		return nil, fmt.Errorf("unmarshal message: %v", err)
	}

	c, err := command.ToModel()
	if err != nil {
		return nil, fmt.Errorf("unmarshal message: %v", err)
	}

	return &c, nil
}

func (c *BrokerConsumer) consumeContext(ctx context.Context) (*kafka.Message, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			msg, err := c.read(ctx)
			if err != nil {
				return nil, err
			}

			if msg.Value == nil {
				continue
			}

			return msg, nil
		}
	}
}
