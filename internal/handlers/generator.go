package handlers

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Generator struct {
	genService     GenerateService
	userService    UserService
	phrasesService PhrasesService
}

func NewGenerator(genService GenerateService, userService UserService, phrasesService PhrasesService) *Generator {
	return &Generator{
		genService:     genService,
		userService:    userService,
		phrasesService: phrasesService,
	}
}

func (g *Generator) GenerateFlower(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	user, err := g.userService.FindByUserName(ctx, update.CallbackQuery.From.Username)
	if err != nil {
		SendError(ctx, b, chatID)
		slog.Error("user not found", update.CallbackQuery.From.Username)
		return
	}

	imgBytes, err := g.genService.GenerateFlower(ctx, user.ID, false)
	if err != nil {
		SendError(ctx, b, chatID)
		slog.Error("generate flower error", err)
		return
	}

	text, err := g.phrasesService.GetRandomText(ctx)
	if err != nil {
		SendError(ctx, b, chatID)
		slog.Error("get random text error", err)
		return
	}

	button, err := g.phrasesService.GetRandomButton(ctx)
	if err != nil {
		SendError(ctx, b, chatID)
		slog.Error("get random button error", err)
		return
	}

	_, err = b.SendPhoto(
		ctx,
		&bot.SendPhotoParams{
			ChatID:  update.CallbackQuery.Message.Message.Chat.ID,
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
