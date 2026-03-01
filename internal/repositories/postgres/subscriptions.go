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

func (s *SubscriptionsRepository) GetReadyUsers(ctx context.Context, time time.Time) ([]entity.User, error) {
	query := `SELECT u.id, u.username, u.chat_id FROM subscriptions s
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
	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Username, &user.ChatID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// TODO: Make move logic in code
func (s *SubscriptionsRepository) MoveTimes(ctx context.Context, subID int) error {
	_, err := s.db.ExecContext(ctx, "UPDATE subscriptions SET last_run = next_run WHERE id = $1", subID)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, "UPDATE subscriptions SET next_run =  WHERE user_id = $1", subID)
	return err
}
