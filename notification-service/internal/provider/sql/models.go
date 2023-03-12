package sql

import domain "notification-service/internal/domain/models"

type NotificationRow struct {
	ID      int64  `db:"id"`
	UserID  string `db:"user_id"`
	Message string `db:"message"`
}

func (r *NotificationRow) ToModel() *domain.Notification {
	return &domain.Notification{
		ID:      r.ID,
		UserID:  r.UserID,
		Message: r.Message,
	}
}
