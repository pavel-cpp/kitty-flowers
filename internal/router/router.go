package router

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pavel-cpp/kitty-flowers/internal/jobs"
	"github.com/pavel-cpp/kitty-flowers/internal/repositories/pollen_ai"
	"github.com/pavel-cpp/kitty-flowers/internal/repositories/postgres"
	"github.com/pavel-cpp/kitty-flowers/internal/services/generator"
	"github.com/pavel-cpp/kitty-flowers/internal/services/phrases"
	"github.com/pavel-cpp/kitty-flowers/internal/services/subscriptions"
	"github.com/pavel-cpp/kitty-flowers/internal/services/user"
	"github.com/robfig/cron/v3"

	"github.com/go-telegram/bot"
	"github.com/pavel-cpp/kitty-flowers/configs"
	"github.com/pavel-cpp/kitty-flowers/internal/handlers"
)

func New(cfg *configs.Config) (*bot.Bot, *cron.Cron) {
	db, err := sql.Open(cfg.DB.Driver, cfg.DB.DBString)
	if err != nil {
		panic(err)
	}

	b, err := bot.New(cfg.TG.Token)
	if err != nil {
		panic(err)
	}

	genRepo := pollen_ai.NewPollerAIRepo(cfg.PollenAI.Token)
	userRepo := postgres.NewUserRepository(db)
	phrasesRepo := postgres.NewPhrasesRepository(db)
	subsRepo := postgres.NewSubscriptionsRepository(db)

	userService := user.New(userRepo)
	genService := generator.New(genRepo, userRepo, cfg.PollenAI.DefaultModel, cfg.PollenAI.GenerateModel)
	phrasesService := phrases.New(phrasesRepo)
	subsService := subscriptions.New(subsRepo)

	userHandler := handlers.NewUser(userService, genService, phrasesService, subsService, cfg.Notifications.NotificationTimes)
	genHandler := handlers.NewGenerator(genService, userService, phrasesService)

	notifyJob := jobs.NewNotifyJob(b, phrasesService, genService, subsService)

	c := cron.New()
	_, err = c.AddJob(cfg.Notifications.CheckFrequency, notifyJob)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, userHandler.Start)

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "generate", bot.MatchTypeExact, genHandler.GenerateFlower)

	return b, c
}
