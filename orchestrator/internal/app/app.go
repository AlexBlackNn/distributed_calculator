package main

import (
	"context"
	"fmt"
	"github.com/levigross/grequests"
	"log/slog"
	"orchestrator/internal/config"
	"orchestrator/internal/services/orchestrator_service"
	"orchestrator/message_broker/rabbit"
	"orchestrator/storage/postgres"
	"orchestrator/storage/redis"
	"os"
	"time"
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

	id, err := application.orchestrationService.CalculationRequest(ctx, "1*1+3")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("1111111", id)

	go application.orchestrationService.ParseResponse(ctx)
	id, err = application.orchestrationService.CalculationRequest(ctx, "1*2+4")
	fmt.Println("22222222", id)

	fmt.Println("_______________________________________________________")
	url := "http://guest:guest@localhost:15672/api/queues/%2f/operation"

	// Send a GET request to the RabbitMQ Management API to get queue details
	resp, err := grequests.Get(url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var data map[string]interface{}
	resp.JSON(&data)

	if consumers, ok := data["consumers"].(float64); ok {
		fmt.Println("Number of consumers connected to the queue:", consumers)
	} else {
		fmt.Println("Consumers information not available for this queue.")
	}
	fmt.Println("_______________________________________________________")
	time.Sleep(10 * time.Second)
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