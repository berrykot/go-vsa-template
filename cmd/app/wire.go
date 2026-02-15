//go:build wireinject
// +build wireinject

package main

import (
	"go-vsa-template/internal/config"
	"go-vsa-template/internal/features/health"
	"go-vsa-template/internal/infrastructure/database"
	"go-vsa-template/internal/infrastructure/logger"
	"go-vsa-template/internal/infrastructure/scheduler"
	"go-vsa-template/internal/infrastructure/server"

	"github.com/google/wire"
)

func InitializeApp() (*App, func(), error) {
	wire.Build(
		config.New,
		logger.New,
		database.NewClient,
		server.New,
		scheduler.New,
		//features below
		health.NewHandler,
		health.NewHandlerProtected,
		health.NewJobHandler,
		NewApp,
	)
	return &App{}, nil, nil
}
