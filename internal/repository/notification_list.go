package repository

import (
	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
)

func (r *NotificationRepository) ListNotifications(filter domain.NotificationListFilter) ([]domain.Notification, error) {
	query := r.db.Model(&database.Notification{})

	if filter.NotifyAtFrom != nil {
		query = query.Where("notify_at >= ?", *filter.NotifyAtFrom)
	}
	if filter.NotifyAtTo != nil {
		query = query.Where("notify_at <= ?", *filter.NotifyAtTo)
	}
	if filter.IsNotified != nil {
		query = query.Where("is_notified = ?", *filter.IsNotified)
	}

	var dbNotifications []database.Notification
	if err := query.Order("notify_at DESC").Find(&dbNotifications).Error; err != nil {
		return nil, err
	}

	notifications := make([]domain.Notification, 0, len(dbNotifications))
	for _, n := range dbNotifications {
		targetUserIDs, err := r.getTargetUserIDs(n.ID)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n.ToDomain(targetUserIDs))
	}

	return notifications, nil
}

func (r *NotificationRepository) getTargetUserIDs(notificationID string) ([]string, error) {
	var targets []database.NotificationTargetUser
	if err := r.db.Where("notification_id = ?", notificationID).Find(&targets).Error; err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(targets))
	for _, t := range targets {
		userIDs = append(userIDs, t.UserID)
	}
	return userIDs, nil
}
