package main

import (
	"context"
	"errors"
	"fmt"
	"go-vsa-template/internal/config"
	"go-vsa-template/internal/features/health"
	"go-vsa-template/internal/infrastructure/scheduler"
	"go-vsa-template/internal/infrastructure/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type App struct {
	Config          *config.Config
	Logger          zerolog.Logger
	Router          *server.Router
	Scheduler       *scheduler.Scheduler
	HealthPublic    *health.Handler
	HealthProtected *health.HandlerProtected
}

func NewApp(cfg *config.Config,
	logger zerolog.Logger,
	router *server.Router,
	scheduler *scheduler.Scheduler,
	// features below
	healthPublic *health.Handler,
	healthProtected *health.HandlerProtected,
	healthJob *health.JobHandler,
) *App {
	// Публичные роуты (без JWT)
	healthPublic.Register(router.Public)

	// Защищённые роуты (Bearer JWT обязателен)
	healthProtected.Register(router.Protected)

	// Scheduler
	healthJob.Register(scheduler)

	return &App{
		Config:          cfg,
		Logger:          logger,
		Router:          router,
		Scheduler:       scheduler,
		HealthPublic:    healthPublic,
		HealthProtected: healthProtected,
	}
}

func (a *App) Run() error {
	// 1. Настраиваем HTTP сервер
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Config.Port),
		Handler: a.Router.Engine,
	}

	// 2. Канал для прослушивания сигналов прерывания (Ctrl+C, Render Shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 3. Запускаем сервер в горутине
	go func() {
		a.Logger.Info().Int("port", a.Config.Port).Msg("server is starting")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Logger.Fatal().Err(err).Msg("failed to listen and serve")
		}
	}()

	a.Scheduler.Start()

	// 4. Блокируемся, пока не прилетит сигнал
	<-quit
	a.Scheduler.Stop()
	a.Logger.Info().Msg("shutting down server...")

	// 5. Даем серверу 5 секунд на завершение текущих задач
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	a.Logger.Info().Msg("server exited gracefully")
	return nil
}
