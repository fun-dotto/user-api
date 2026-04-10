package repository

import (
	"errors"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
	"gorm.io/gorm"
)

func (r *NotificationRepository) UpdateNotification(notification domain.Notification) (domain.Notification, error) {
	tx := r.db.Begin()

	var existing database.Notification
	if err := tx.First(&existing, "id = ?", notification.ID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Notification{}, domain.ErrNotFound
		}
		return domain.Notification{}, err
	}

	dbNotification := database.NotificationFromDomain(notification)
	if err := tx.Save(&dbNotification).Error; err != nil {
		tx.Rollback()
		return domain.Notification{}, err
	}

	if err := tx.Where("notification_id = ?", notification.ID).Delete(&database.NotificationTargetUser{}).Error; err != nil {
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
