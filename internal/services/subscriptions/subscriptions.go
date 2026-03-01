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
	GetReadyUsers(ctx context.Context, time time.Time) ([]entity.User, error)
	UpdateSubscription(ctx context.Context, userID uuid.UUID, timestamp time.Time) error
}

type Service struct {
	repo SubscriptionsRepository
}

func New(repo SubscriptionsRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Subscribe(ctx context.Context, userID uuid.UUID, timestamp time.Time) error {
	return s.repo.CreateSubscription(ctx, userID, timestamp)
}

func (s *Service) NotifyUsers(ctx context.Context, notifyFunc func(context.Context, entity.User)) error {
	users, err := s.repo.GetReadyUsers(ctx, time.Now())
	fmt.Println(users)
	if err != nil {
		return err
	}
	for _, user := range users {
		notifyFunc(ctx, user)
		err = s.repo.UpdateSubscription(ctx, user.ID, time.Now())
		if err != nil {
			slog.Error("failed to update subscription", "user_id", user.ID, "error", err)
		}
	}
	return nil
}
