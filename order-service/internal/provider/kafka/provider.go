package kafka

import (
	"broker/consumer"
	"fmt"

	"time"
)

const groupId = "eventGroupStore"

type kafkaProvider struct {
	topic    string
	cfg      *broker.NetworkConfig
	consumer consumer.Consumer
}

func NewKafkaProvider(cfg *broker.NetworkConfig, topic string) *kafkaProvider {
	return &kafkaProvider{
		cfg:   cfg,
		topic: topic,
	}
}

func (e *kafkaProvider) StartConsume(ctx context.Context) (<-chan dto.FacadeRequestPayload, <-chan error, error) {
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

func (e *kafkaProvider) consume(ctx context.Context) (<-chan dto.FacadeRequestPayload, <-chan error) {
	payloadCh := make(chan dto.FacadeRequestPayload)
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
				start := time.Now()
				if err != nil {
					errCh <- fmt.Errorf("consume message: %v", err)
					continue
				}

				var batch []LogMessage

				if err = json.Unmarshal(message.Value, &batch); err != nil {
					errCh <- fmt.Errorf("unmarshal message: %v", err)
					continue
				}

				log.Debugf("Log message (batch) came: %d events", len(batch))

				events := make([]*dto.EventMessageDTO, 0, len(batch))

				for _, v := range batch {
					event, err := v.ToEventMessageFormat()
					if err != nil {
						errCh <- fmt.Errorf("convert message: %v", err)
						continue
					}

					if event.UserID != "" {
						events = append(events, event)
					}
				}

				metric := dto.FacadeRequestPayload{
					Payload: events,
					Start:   start,
				}

				payloadCh <- metric
			}
		}
	}()

	return payloadCh, errCh
}
