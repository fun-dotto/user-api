package repository

import (
	"context"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
	"gorm.io/gorm"
)

func (r *NotificationRepository) DispatchNotifications(ctx context.Context, ids []string) ([]domain.Notification, error) {
	uniqueIDs := uniqueStrings(ids)
	if len(uniqueIDs) == 0 {
		return []domain.Notification{}, nil
	}

	var dbNotifications []database.Notification

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", uniqueIDs).Find(&dbNotifications).Error; err != nil {
			return err
		}
		if len(dbNotifications) == 0 {
			return nil
		}

		existingIDs := make([]string, 0, len(dbNotifications))
		for _, n := range dbNotifications {
			existingIDs = append(existingIDs, n.ID)
		}

		if err := tx.Model(&database.Notification{}).
			Where("id IN ?", existingIDs).
			Update("is_notified", true).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(dbNotifications) == 0 {
		return []domain.Notification{}, nil
	}

	existingIDs := make([]string, 0, len(dbNotifications))
	for _, n := range dbNotifications {
		existingIDs = append(existingIDs, n.ID)
	}

	var allTargets []database.NotificationTargetUser
	if err := r.db.WithContext(ctx).Where("notification_id IN ?", existingIDs).Find(&allTargets).Error; err != nil {
		return nil, err
	}

	targetMap := make(map[string][]string)
	for _, t := range allTargets {
		targetMap[t.NotificationID] = append(targetMap[t.NotificationID], t.UserID)
	}

	notifications := make([]domain.Notification, 0, len(dbNotifications))
	for _, n := range dbNotifications {
		n.IsNotified = true
		notifications = append(notifications, n.ToDomain(targetMap[n.ID]))
	}

	return notifications, nil
}
