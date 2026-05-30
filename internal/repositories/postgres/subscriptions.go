package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/pavel-cpp/kitty-flowers/internal/entity"
)

var (
	ErrSubscriptionNotCreated = errors.New("subscription not created")
)

type SubscriptionsRepository struct {
	db *sql.DB
}

func NewSubscriptionsRepository(db *sql.DB) *SubscriptionsRepository {
	return &SubscriptionsRepository{db: db}
}

func (s *SubscriptionsRepository) IsSubscribed(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM subscriptions WHERE user_id = $1`
	row := s.db.QueryRowContext(ctx, query, userID)
	if row.Err() != nil {
		return false, row.Err()
	}
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *SubscriptionsRepository) CreateSubscription(ctx context.Context, userID uuid.UUID, timestamp time.Time) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO subscriptions (user_id, next_run) VALUES ($1, $2)", userID, timestamp)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			panic(err)
		}
		slog.Error("subscription not created", err)
		return ErrSubscriptionNotCreated
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	return nil
}

func (s *SubscriptionsRepository) GetUnnotifiedUsers(ctx context.Context, time time.Time) ([]entity.UserNotification, error) {
	query := `SELECT u.id, u.username, u.chat_id, s.id, s.next_run FROM subscriptions s
			LEFT JOIN users u ON u.id = s.user_id
			WHERE s.active IS TRUE AND s.next_run < $1`
	rows, err := s.db.QueryContext(ctx, query, time)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)
	var users []entity.UserNotification
	for rows.Next() {
		var user entity.UserNotification
		err = rows.Scan(&user.ID, &user.Username, &user.ChatID, &user.NotificationID, &user.CurrentRun)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *SubscriptionsRepository) UpdateUserNotificationTime(ctx context.Context, notificationID int, nextRun time.Time) error {
	_, err := s.db.ExecContext(ctx, "UPDATE subscriptions SET last_run = next_run WHERE id = $1", notificationID)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, "UPDATE subscriptions SET next_run = $1 WHERE id = $2", nextRun, notificationID)
	return err
}
