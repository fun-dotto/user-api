package service

import (
	"context"

	"github.com/fun-dotto/user-api/internal/domain"
)

func (s *NotificationService) DispatchNotifications(ctx context.Context, ids []string) ([]domain.Notification, error) {
	return s.repo.DispatchNotifications(ctx, ids)
}
