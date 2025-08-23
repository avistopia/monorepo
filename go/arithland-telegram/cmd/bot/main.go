package main

import (
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/models"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/clients"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/core"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/handler"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Panic("failed to load config")
	}

	var connection gorm.Dialector

	switch config.Database.Type {
	case "sqlite":
		connection = sqlite.Open(config.Database.DSN)
	case "postgres":
		connection = postgres.Open(config.Database.DSN)
	default:
		logrus.Panicf("unknown database type %q", config.Database.Type)
	}

	db, err := gorm.Open(connection, &gorm.Config{})
	if err != nil {
		logrus.WithError(err).Panic("failed to connect to database")
	}

	userRepo, err := models.NewUserRepo(db)
	if err != nil {
		logrus.WithError(err).Panic("failed to initialize user repo")
	}

	questionRepo, err := models.NewQuestionRepo(db)
	if err != nil {
		logrus.WithError(err).Panic("failed to initialize question repo")
	}

	userQuestionRepo, err := models.NewUserQuestionRepo(db)
	if err != nil {
		logrus.WithError(err).Panic("failed to initialize user repo")
	}

	telegram, err := clients.NewTelegram(config.Telegram.Token)
	if err != nil {
		logrus.WithError(err).Panic("failed to initialize telegram client")
	}

	adminUsernames := map[string]struct{}{}
	for _, u := range config.Telegram.AdminUsernames {
		adminUsernames[u] = struct{}{}
	}

	flow, err := core.NewService(
		telegram, adminUsernames, config.Telegram.DefaultPhotoID, db, userRepo, questionRepo, userQuestionRepo,
	).Flow()
	if err != nil {
		logrus.WithError(err).Panic("failed to get the core service flows")
	}

	handler.NewHandler(telegram, userRepo, flow).Listen()
}
