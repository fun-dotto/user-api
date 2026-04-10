package repository

import (
	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
	"github.com/google/uuid"
)

func (r *NotificationRepository) CreateNotification(notification domain.Notification) (domain.Notification, error) {
	notification.ID = uuid.New().String()

	dbNotification := database.NotificationFromDomain(notification)

	tx := r.db.Begin()

	if err := tx.Create(&dbNotification).Error; err != nil {
		tx.Rollback()
		return domain.Notification{}, err
	}

	for _, userID := range notification.TargetUserIDs {
		target := database.NotificationTargetUser{
			NotificationID: notification.ID,
			UserID:         userID,
		}
		if err := tx.Create(&target).Error; err != nil {
			tx.Rollback()
			return domain.Notification{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.Notification{}, err
	}

	return dbNotification.ToDomain(notification.TargetUserIDs), nil
}
