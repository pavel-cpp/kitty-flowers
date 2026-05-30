package subscriptions

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/pavel-cpp/kitty-flowers/internal/entity"
)

type SubscriptionsRepository interface {
	CreateSubscription(ctx context.Context, userID uuid.UUID, timestamp time.Time) error
	GetUnnotifiedUsers(ctx context.Context, time time.Time) ([]entity.UserNotification, error)
	UpdateUserNotificationTime(ctx context.Context, notificationID int, nextRun time.Time) error
	IsSubscribed(ctx context.Context, userID uuid.UUID) (bool, error)
}

type Service struct {
	repo SubscriptionsRepository
}

func New(repo SubscriptionsRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) IsSubscribed(ctx context.Context, userID uuid.UUID) (bool, error) {
	return s.repo.IsSubscribed(ctx, userID)
}

func (s *Service) Subscribe(ctx context.Context, userID uuid.UUID, timeOfDay time.Time) error {
	now := time.Now()
	year, month, day := now.Date()
	timestamp := timeOfDay.AddDate(year, int(month)-1, day-1)
	if now.After(timestamp) {
		timestamp = timestamp.AddDate(0, 0, 1)
	}
	return s.repo.CreateSubscription(ctx, userID, timestamp)
}

func (s *Service) NotifyUsers(ctx context.Context, notifyFunc func(context.Context, entity.User)) error {
	users, err := s.repo.GetUnnotifiedUsers(ctx, time.Now())
	fmt.Println(users)
	if err != nil {
		return err
	}
	for _, user := range users {
		notifyFunc(ctx, user.User)
		err = s.repo.UpdateUserNotificationTime(ctx, user.NotificationID, user.CurrentRun.AddDate(0, 0, 1))
		if err != nil {
			slog.Error("failed to update subscription", "user_id", user.ID, "error", err)
		}
	}
	return nil
}
