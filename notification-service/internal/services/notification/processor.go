package notification

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	domain "notification-service/internal/domain/models"
	"notification-service/internal/provider"
	"sync"
)

type Config struct {
	StorageProv         provider.StorageProvider
	CommandConsumerProv provider.BrokerConsumerProvider
	UserConsumerProv    provider.BrokerConsumerProvider
}

type Processor struct {
	startOnce           sync.Once
	commandCh           <-chan domain.Order
	userUpdateCh        <-chan domain.User
	storageProv         provider.StorageProvider
	commandConsumerProv provider.BrokerConsumerProvider
	userConsumerProv    provider.BrokerConsumerProvider
}

func NewProcessor(cfg Config) *Processor {
	return &Processor{
		storageProv:         cfg.StorageProv,
		commandConsumerProv: cfg.CommandConsumerProv,
		userConsumerProv:    cfg.UserConsumerProv,
	}
}

func (p *Processor) Run(ctx context.Context) {
	errCh := make(chan error, 0)
	payloadCh, err := p.commandConsumerProv.StartConsume(ctx, errCh)
	userUpdateCh, err := p.userConsumerProv.StartConsumeUserUpdate(ctx, errCh)
	if err != nil {
		log.Errorf("failed consumer: %v", err)
	}

	p.commandCh = payloadCh
	p.userUpdateCh = userUpdateCh

	p.startOnce.Do(func() {
		go p.start(ctx)
	})
	fmt.Println(fmt.Sprintf("processor started"))
}

func (p *Processor) start(ctx context.Context) func() {
	for {
		select {
		case command := <-p.commandCh:
			p.executeCommand(ctx, command)
		case user := <-p.userUpdateCh:
			p.updateUser(ctx, user)
		case <-ctx.Done():
			log.Infof("Contex faired! Stopping processor service...")
			break
		}
	}
}

func (p *Processor) executeCommand(ctx context.Context, command domain.Order) {
	switch command.Status {
	case domain.Success, domain.Canceled:
		user, err := p.storageProv.GetUserByID(ctx, command.UserID)
		if err != nil {
			log.Errorf("execute command: %s, with orrder: %v", err, command)
		}
		if user == nil {
			log.Errorf("execute command: %s, with orrder: %v", err, command)
		}

		message := fmt.Sprintf("Order %d %s", command.ID, command.Status)

		p.storageProv.Create(ctx, domain.Notification{
			Mail:    user.Mail,
			Message: message,
		})
	default:
	}
}

func (p *Processor) updateUser(ctx context.Context, user domain.User) {
	if err := p.storageProv.UpdateUserInfo(ctx, user); err != nil {
		log.Errorf("failed update user: %s", err)
	}
}
