package generator

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/google/uuid"
)

type (
	UserRepository interface {
		IncrementDelivery(ctx context.Context, id uuid.UUID) error
	}

	GeneratorRepo interface {
		GetFlower(ctx context.Context, hexColor string, model string) ([]byte, error)
	}
)

type Service struct {
	genRepo       GeneratorRepo
	userRepo      UserRepository
	defaultModel  string
	generateModel string
}

func New(genRepo GeneratorRepo, userRepo UserRepository, defaultModel, generateModel string) *Service {
	return &Service{
		genRepo:       genRepo,
		userRepo:      userRepo,
		defaultModel:  defaultModel,
		generateModel: generateModel,
	}
}

func (s *Service) GenerateFlower(ctx context.Context, userID uuid.UUID, initial bool) ([]byte, error) {
	var img []byte
	model := s.defaultModel
	if initial {
		model = s.generateModel
	}

	img, err := s.genRepo.GetFlower(ctx, fmt.Sprintf("%x", rand.Uint32()), model)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.IncrementDelivery(ctx, userID)
	if err != nil {
		slog.Warn("delivery not incremented", userID, err.Error())
	}

	return img, nil
}
