package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type PhrasesRepository struct {
	db *sql.DB
}

func NewPhrasesRepository(db *sql.DB) *PhrasesRepository {
	return &PhrasesRepository{db: db}
}

func (pr *PhrasesRepository) prasesAmount(ctx context.Context, phrasesType string) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) from %s", phrasesType)

	row := pr.db.QueryRowContext(ctx, query)
	if err := row.Err(); err != nil {
		return 0, err
	}
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (pr *PhrasesRepository) getPrase(ctx context.Context, id int, phraseType string) (string, error) {
	query := fmt.Sprintf("SELECT phrase FROM %s WHERE id = $1", phraseType)

	row := pr.db.QueryRowContext(ctx, query, id)
	if err := row.Err(); err != nil {
		return "", err
	}
	var text string
	if err := row.Scan(&text); err != nil {
		return "", err
	}
	return text, nil
}

func (pr *PhrasesRepository) TextPhrasesAmount(ctx context.Context) (int, error) {
	return pr.prasesAmount(ctx, "text_phrases")
}

func (pr *PhrasesRepository) GetTextPhrase(ctx context.Context, id int) (string, error) {
	return pr.getPrase(ctx, id, "text_phrases")
}

func (pr *PhrasesRepository) ButtonPhrasesAmount(ctx context.Context) (int, error) {
	return pr.prasesAmount(ctx, "button_phrases")
}

func (pr *PhrasesRepository) GetButtonPhrase(ctx context.Context, id int) (string, error) {
	return pr.getPrase(ctx, id, "button_phrases")
}
