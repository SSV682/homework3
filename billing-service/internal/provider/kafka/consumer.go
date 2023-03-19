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

const (
	billingServiceName = "billing_service"
	orderServiceName   = "order_service"
)

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

func (c *BrokerConsumer) read(ctx context.Context) ([]byte, []kafka.Header, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, nil, err
	}
	return msg.Value, msg.Headers, nil
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
				message, headers, err := c.consumeContext(ctx)
				if err != nil {
					fmt.Println(err)
					errCh <- fmt.Errorf("consume message: %v", err)
					continue
				}

				var messageType string
				for _, v := range headers {
					if v.Key == "publisher" {
						messageType = string(v.Value)
					}
				}

				switch messageType {
				case billingServiceName:
					println("yeee order_service")
					command, err := getUserCommand(message)
					if err != nil {
						println(err)
						continue
					}
					println(fmt.Sprintf("command: %#v", command))
					payloadUserCh <- command
				case orderServiceName:
					println("yeee order_service")
					command, err := getBillingCommand(message)
					if err != nil {
						continue
					}
					println(fmt.Sprintf("command: %#v", command))
					payloadCh <- command
				default:
				}
			}
		}
	}()

	return payloadCh, payloadUserCh, errCh, nil
}

func getBillingCommand(message []byte) (domain.RequestCommand, error) {
	var command BillingRequestCommand

	if err := json.Unmarshal(message, &command); err != nil {
		return domain.RequestCommand{}, fmt.Errorf("unmarshal message: %v", err)
	}

	c, err := command.ToModel()
	if err != nil {
		return domain.RequestCommand{}, fmt.Errorf("unmarshal message: %v", err)
	}

	return c, nil
}

func getUserCommand(message []byte) (domain.Account, error) {
	var command UserRequestCommand

	if err := json.Unmarshal(message, &command); err != nil {
		fmt.Printf("%s", err)
		return domain.Account{}, fmt.Errorf("unmarshal message: %v", err)
	}

	c, err := command.ToModel()
	if err != nil {
		return domain.Account{}, fmt.Errorf("unmarshal message: %v", err)
	}

	return c, nil
}

func (c *BrokerConsumer) consumeContext(ctx context.Context) ([]byte, []kafka.Header, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
			msg, headers, err := c.read(ctx)
			if err != nil {
				return nil, nil, err
			}
			fmt.Println(msg)
			fmt.Println(headers)
			if msg == nil {
				continue
			}

			return msg, headers, nil
		}
	}
}
