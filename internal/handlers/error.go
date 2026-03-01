package handlers

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
)

func SendError(ctx context.Context, b *bot.Bot, chatID int64) {
	_, err := b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Кажется у нас какие-то неполадки(",
		},
	)
	if err != nil {
		slog.Error("send message error", err)
	}
}
