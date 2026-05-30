package handlers

import (
	"bytes"
	"context"
	"log/slog"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
	"github.com/pavel-cpp/kitty-flowers/internal/entity"
)

type (
	UserService interface {
		RegisterUser(ctx context.Context, user entity.User) (uuid.UUID, error)
		FindByUserName(ctx context.Context, username string) (entity.User, error)
	}

	GenerateService interface {
		GenerateFlower(ctx context.Context, userID uuid.UUID, initial bool) ([]byte, error)
	}

	PhrasesService interface {
		GetRandomText(ctx context.Context) (string, error)
		GetRandomButton(ctx context.Context) (string, error)
	}

	SubscriptionService interface {
		Subscribe(ctx context.Context, userID uuid.UUID, timestamp time.Time) error
	}
)

type User struct {
	userService      UserService
	genService       GenerateService
	phrasesService   PhrasesService
	subsService      SubscriptionService
	defaultSubsTimes []time.Time
}

func NewUser(userService UserService, genService GenerateService, phrasesService PhrasesService, subsService SubscriptionService, defaultSubsTimes []time.Time) *User {
	return &User{
		userService:      userService,
		genService:       genService,
		phrasesService:   phrasesService,
		subsService:      subsService,
		defaultSubsTimes: defaultSubsTimes,
	}
}

func (u *User) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	id, err := u.userService.RegisterUser(
		ctx,
		entity.User{Username: update.Message.From.Username, ChatID: int(update.Message.Chat.ID)},
	)
	if err != nil {
		slog.Error("register user error", err)
		return
	}

	for _, t := range u.defaultSubsTimes {
		err := u.subsService.Subscribe(ctx, id, t)
		if err != nil {
			slog.Error("subscribe error", err)
		}
	}

	imgBytes, err := u.genService.GenerateFlower(ctx, id, true)
	if err != nil {
		SendError(ctx, b, update.Message.Chat.ID)
		slog.Error("generate flower error", err)
		return
	}

	text, err := u.phrasesService.GetRandomText(ctx)
	if err != nil {
		SendError(ctx, b, update.Message.Chat.ID)
		slog.Error("get random text error", err)
		return
	}

	button, err := u.phrasesService.GetRandomButton(ctx)
	if err != nil {
		SendError(ctx, b, update.Message.Chat.ID)
		slog.Error("get random button error", err)
		return
	}
	_, err = b.SendPhoto(
		ctx,
		&bot.SendPhotoParams{
			ChatID:  update.Message.Chat.ID,
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
