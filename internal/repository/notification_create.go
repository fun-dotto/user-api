package repository

import (
	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *NotificationRepository) CreateNotification(notification domain.Notification) (domain.Notification, error) {
	notification.ID = uuid.New().String()

	dbNotification := database.NotificationFromDomain(notification)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dbNotification).Error; err != nil {
			return err
		}

		for _, userID := range notification.TargetUserIDs {
			target := database.NotificationTargetUser{
				NotificationID: notification.ID,
				UserID:         userID,
			}
			if err := tx.Create(&target).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return domain.Notification{}, err
	}

	return dbNotification.ToDomain(notification.TargetUserIDs), nil
}
