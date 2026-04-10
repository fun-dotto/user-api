package repository

import (
	"errors"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
	"gorm.io/gorm"
)

func (r *NotificationRepository) DeleteNotification(id string) error {
	tx := r.db.Begin()

	var existing database.Notification
	if err := tx.First(&existing, "id = ?", id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	if err := tx.Where("notification_id = ?", id).Delete(&database.NotificationTargetUser{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&existing).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
