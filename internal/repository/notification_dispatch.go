package repository

import (
	"context"
	"strings"
	"time"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
)

func (r *NotificationRepository) GetNotificationsByIDs(ctx context.Context, ids []string) ([]domain.Notification, error) {
	uniqueIDs := uniqueStrings(ids)
	if len(uniqueIDs) == 0 {
		return []domain.Notification{}, nil
	}

	var dbNotifications []database.Notification
	if err := r.db.WithContext(ctx).Where("id IN ?", uniqueIDs).Find(&dbNotifications).Error; err != nil {
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

	targetMap := make(map[string][]domain.NotificationTargetUser)
	for _, t := range allTargets {
		targetMap[t.NotificationID] = append(targetMap[t.NotificationID], domain.NotificationTargetUser{
			UserID:     t.UserID,
			NotifiedAt: t.NotifiedAt,
		})
	}

	notifications := make([]domain.Notification, 0, len(dbNotifications))
	for _, n := range dbNotifications {
		notifications = append(notifications, n.ToDomain(targetMap[n.ID]))
	}

	return notifications, nil
}

func (r *NotificationRepository) DispatchNotifications(ctx context.Context, deliveries map[string][]string) ([]domain.Notification, error) {
	if len(deliveries) == 0 {
		return []domain.Notification{}, nil
	}

	now := time.Now()
	placeholders := make([]string, 0)
	args := make([]any, 0)
	args = append(args, now)
	for nid, userIDs := range deliveries {
		uniqueUsers := uniqueStrings(userIDs)
		for _, uid := range uniqueUsers {
			placeholders = append(placeholders, "(?, ?)")
			args = append(args, nid, uid)
		}
	}
	if len(placeholders) == 0 {
		return []domain.Notification{}, nil
	}

	sql := "UPDATE notification_target_users SET notified_at = ? " +
		"WHERE (notification_id, user_id) IN (" + strings.Join(placeholders, ", ") + ") " +
		"AND notified_at IS NULL RETURNING notification_id"

	var updated []database.NotificationTargetUser
	if err := r.db.WithContext(ctx).Raw(sql, args...).Scan(&updated).Error; err != nil {
		return nil, err
	}

	if len(updated) == 0 {
		return []domain.Notification{}, nil
	}

	notificationIDs := make([]string, 0, len(updated))
	seen := make(map[string]struct{}, len(updated))
	for _, row := range updated {
		if _, ok := seen[row.NotificationID]; ok {
			continue
		}
		seen[row.NotificationID] = struct{}{}
		notificationIDs = append(notificationIDs, row.NotificationID)
	}

	return r.GetNotificationsByIDs(ctx, notificationIDs)
}
