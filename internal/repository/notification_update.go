package repository

import (
	"errors"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
	"gorm.io/gorm"
)

func (r *NotificationRepository) UpdateNotification(notification domain.Notification) (domain.Notification, error) {
	var dbNotification database.Notification

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing database.Notification
		if err := tx.First(&existing, "id = ?", notification.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrNotFound
			}
			return err
		}

		dbNotification = database.NotificationFromDomain(notification)

		if err := tx.Omit("IsNotified").Save(&dbNotification).Error; err != nil {
			return err
		}

		if err := tx.Where("notification_id = ?", notification.ID).Delete(&database.NotificationTargetUser{}).Error; err != nil {
			return err
		}

		for _, userID := range uniqueStrings(notification.TargetUserIDs) {
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
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Notification{}, domain.ErrNotFound
		}
		return domain.Notification{}, err
	}

	return dbNotification.ToDomain(notification.TargetUserIDs), nil
}
