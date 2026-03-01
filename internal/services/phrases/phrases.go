package phrases

import (
	"context"
	"math/rand"
)

type PhraseRepository interface {
	TextPhrasesAmount(ctx context.Context) (int, error)
	GetTextPhrase(ctx context.Context, id int) (string, error)
	ButtonPhrasesAmount(ctx context.Context) (int, error)
	GetButtonPhrase(ctx context.Context, id int) (string, error)
}

type Service struct {
	repo PhraseRepository
}

func New(repo PhraseRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetRandomText(ctx context.Context) (string, error) {
	amount, err := s.repo.TextPhrasesAmount(ctx)
	if err != nil {
		return "", err
	}
	return s.repo.GetTextPhrase(ctx, rand.Int()%amount)
}

func (s *Service) GetRandomButton(ctx context.Context) (string, error) {
	amount, err := s.repo.ButtonPhrasesAmount(ctx)
	if err != nil {
		return "", err
	}
	return s.repo.GetButtonPhrase(ctx, rand.Int()%amount)
}
