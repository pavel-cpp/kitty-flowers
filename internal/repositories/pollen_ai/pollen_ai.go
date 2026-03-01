package pollen_ai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type PollerAI struct {
	token string
}

func NewPollerAIRepo(token string) *PollerAI {
	return &PollerAI{
		token: token,
	}
}

func (p *PollerAI) GetFlower(ctx context.Context, hexColor string, model string) ([]byte, error) {
	url := fmt.Sprintf("https://gen.pollinations.ai/image/unique%%20flower%%20with%%20color%%20%s?model=%s", hexColor, model)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.token))
	req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		slog.Warn("Error fetching pollen ai")
		return nil, errors.New(res.Status)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	img, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}
