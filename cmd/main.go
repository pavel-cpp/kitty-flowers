package main

import (
	"context"

	"github.com/pavel-cpp/kitty-flowers/configs"
	"github.com/pavel-cpp/kitty-flowers/internal/router"
)

func main() {
	cfg := configs.LoadConfig()
	bot, scheduler := router.New(&cfg)
	scheduler.Start()
	bot.Start(context.Background())
}
