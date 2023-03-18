package user

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"user-service/internal/domain/models"
	"user-service/internal/provider"
)

type ServiceConfig struct {
	SqlProv               provider.SqlUserProvider
	BrokerProv            provider.BrokerProducerProvider
	BillingTopicName      string
	NotificationTopicName string
}

type userService struct {
	sqlProv               provider.SqlUserProvider
	brokerProv            provider.BrokerProducerProvider
	billingTopicName      string
	notificationTopicName string
}

func NewUserService(cfg ServiceConfig) *userService {
	return &userService{
		sqlProv:               cfg.SqlProv,
		brokerProv:            cfg.BrokerProv,
		billingTopicName:      cfg.BillingTopicName,
		notificationTopicName: cfg.NotificationTopicName,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) (string, error) {
	password := user.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("password hashing failed:  %s", err)
	}

	user.Password = string(hashedPassword)

	log.Debugf("user: %#v", user)

	i, err := s.sqlProv.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	user.ID = i

	err = s.brokerProv.SendCommand(ctx, *user, []string{s.billingTopicName, s.notificationTopicName})
	if err != nil {
		return "", fmt.Errorf("couldnt create account: %s", err)
	}

	return i, nil

}

func (s *userService) GetUser(ctx context.Context, userID string) (models.User, error) {
	user, err := s.sqlProv.GetUser(ctx, userID)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by id %s: %v", userID, err)
	}
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	if err := s.sqlProv.DeleteUser(ctx, userID); err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUser(ctx context.Context, userID string, user *models.User) error {
	password := user.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hashing failed:  %s", err)
	}

	user.Password = string(hashedPassword)

	err = s.brokerProv.SendCommand(ctx, *user, []string{s.notificationTopicName})
	if err != nil {
		log.Errorf("couldnt create account: %s", err)
	}

	if err := s.sqlProv.UpdateUser(ctx, userID, user); err != nil {
		return err
	}
	return nil
}
