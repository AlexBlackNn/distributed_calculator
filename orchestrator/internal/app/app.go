package main

import (
	"context"
	"fmt"
	"log/slog"
	"orchestrator/internal/config"
	"orchestrator/internal/services/orchestrator_service"
	"orchestrator/message_broker/rabbit"
	"orchestrator/storage/postgres"
	"orchestrator/storage/redis"
	"os"
)

type App struct {
	orchestrationService *orchestrator_service.OrchestratorService
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	//init storage
	operationSettingsStorage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}
	//init cache
	operationCache := redis.New(cfg)

	messageBroker, err := rabbit.New("test")
	if err != nil {
		panic(err)
	}

	//init orchestrator_service orchestrator_service (orchestrator_service)
	orchestrationService := orchestrator_service.New(
		log,
		operationSettingsStorage,
		operationSettingsStorage,
		operationCache,
		messageBroker,
		cfg,
	)

	return &App{
		orchestrationService: orchestrationService,
	}
}

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()
	// init logger
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	application := New(log, cfg)

	id, err := application.orchestrationService.CalculationRequest(ctx, "1*1+(2*2)+3")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id)
}

const (
	envLocal = "local"
	envDemo  = "demo"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{
					Level:     slog.LevelDebug,
					AddSource: true,
				},
			),
		)
	case envDemo:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{
					Level:     slog.LevelDebug,
					AddSource: true,
				},
			),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{
					Level:     slog.LevelInfo,
					AddSource: true,
				},
			),
		)
	}
	return log
}
