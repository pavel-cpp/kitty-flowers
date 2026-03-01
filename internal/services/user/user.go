package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/pavel-cpp/kitty-flowers/internal/entity"
)

const notFound = "not found"

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (uuid.UUID, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	FindByUserName(ctx context.Context, username string) (entity.User, error)
}

type Service struct {
	repo UserRepository
}

func New(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterUser(ctx context.Context, user entity.User) (uuid.UUID, error) {
	foundUser, err := s.repo.FindByUserName(ctx, user.Username)
	if err != nil && err.Error() == notFound {
		return s.repo.Create(ctx, user)
	}
	return foundUser.ID, err
}

func (s *Service) FindByUserName(ctx context.Context, username string) (entity.User, error) {
	user, err := s.repo.FindByUserName(ctx, username)
	if err != nil && err.Error() == notFound {
		return entity.User{}, errors.New(notFound)
	}
	return user, err
}

func (s *Service) FindByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil && err.Error() == notFound {
		return entity.User{}, errors.New(notFound)
	}
	return user, err
}
