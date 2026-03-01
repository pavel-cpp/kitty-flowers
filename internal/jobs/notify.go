package jobs

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
	"github.com/pavel-cpp/kitty-flowers/internal/entity"
)

type (
	GenerateService interface {
		GenerateFlower(ctx context.Context, userID uuid.UUID, initial bool) ([]byte, error)
	}

	PhrasesService interface {
		GetRandomText(ctx context.Context) (string, error)
		GetRandomButton(ctx context.Context) (string, error)
	}

	SubscriptionsService interface {
		NotifyUsers(ctx context.Context, notifyFunc func(context.Context, entity.User)) error
	}
)

type NotifyJob struct {
	b              *bot.Bot
	phrasesService PhrasesService
	genService     GenerateService
	subsService    SubscriptionsService
}

func NewNotifyJob(b *bot.Bot, phrasesService PhrasesService, genService GenerateService, subsService SubscriptionsService) *NotifyJob {
	return &NotifyJob{
		b:              b,
		phrasesService: phrasesService,
		genService:     genService,
		subsService:    subsService,
	}
}

func (j *NotifyJob) NotifyUser(ctx context.Context, user entity.User) {
	imgBytes, err := j.genService.GenerateFlower(ctx, user.ID, false)
	if err != nil {
		slog.Error("generate flower error", err)
		return
	}

	text, err := j.phrasesService.GetRandomText(ctx)
	if err != nil {
		slog.Error("get random text error", err)
		return
	}

	button, err := j.phrasesService.GetRandomButton(ctx)
	if err != nil {
		slog.Error("get random button error", err)
		return
	}

	_, err = j.b.SendPhoto(
		ctx,
		&bot.SendPhotoParams{
			ChatID:  user.ChatID,
			Photo:   &models.InputFileUpload{Filename: "image.png", Data: bytes.NewReader(imgBytes)},
			Caption: text,
			ReplyMarkup: &models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					{{
						Text:         button,
						CallbackData: "generate",
					}},
				},
			},
		},
	)
	if err != nil {
		slog.Error("send message error", err)
	}
}

func (j *NotifyJob) Run() {
	fmt.Println("job started")
	err := j.subsService.NotifyUsers(
		context.Background(),
		j.NotifyUser,
	)
	if err != nil {
		slog.Error("notify users error", err)
		return
	}
}
