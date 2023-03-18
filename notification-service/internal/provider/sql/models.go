package sql

import domain "notification-service/internal/domain/models"

type NotificationRow struct {
	ID      int64  `db:"id"`
	Mail    string `db:"mail"`
	Message string `db:"message"`
}

func (r *NotificationRow) ToModel() *domain.Notification {
	return &domain.Notification{
		ID:      r.ID,
		Mail:    r.Mail,
		Message: r.Message,
	}
}

type UserInfoRow struct {
	ID   string `db:"user_id"`
	Mail string `db:"mail"`
}

func (r *UserInfoRow) ToModel() *domain.User {
	return &domain.User{
		ID:   r.ID,
		Mail: r.Mail,
	}
}
