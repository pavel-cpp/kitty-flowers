package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/pavel-cpp/kitty-flowers/internal/entity"
)

var ErrNotFound = errors.New("not found")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(ctx context.Context, user entity.User) (uuid.UUID, error) {
	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}
	row := tx.QueryRowContext(ctx, "INSERT INTO users (username, chat_id) VALUES ($1, $2) RETURNING id;", user.Username, user.ChatID)
	if row.Err() != nil {
		err = tx.Rollback()
		if err != nil {
			return uuid.Nil, err
		}
		return uuid.Nil, err
	}
	var id uuid.UUID
	if err = row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	err = tx.Commit()
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (ur *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	row := ur.db.QueryRowContext(ctx, "SELECT username, chat_id FROM users WHERE id = $1;", id)
	if row.Err() != nil {
		return entity.User{}, row.Err()
	}
	user := entity.User{
		ID: id,
	}
	err := row.Scan(&user.Username, &user.ChatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, ErrNotFound
		}
		return entity.User{}, err
	}
	return user, nil
}

func (ur *UserRepository) FindByUserName(ctx context.Context, username string) (entity.User, error) {
	row := ur.db.QueryRowContext(ctx, "SELECT id, username, chat_id FROM users WHERE username = $1;", username)
	if row.Err() != nil {
		return entity.User{}, row.Err()
	}
	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.ChatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, ErrNotFound
		}
		return entity.User{}, err
	}
	return user, nil
}

func (ur *UserRepository) IncrementDelivery(ctx context.Context, id uuid.UUID) error {
	_, err := ur.db.Exec("UPDATE user_stats SET images_delivered = images_delivered + 1 where user_id = $1;", id)
	return err
}
