package app

import (
	"log/slog"
	"orchestrator/internal/config"
	"orchestrator/internal/services/monitoring_service"
	"orchestrator/internal/services/orchestrator_service"
	"orchestrator/internal/services/settings_service"
	"orchestrator/message_broker/rabbit"
	"orchestrator/storage/postgres"
	"os"
)

type App struct {
	OrchestrationService *orchestrator_service.OrchestratorService
	MonitoringService    *monitoring_service.MonitoringService
	SettingService       *settings_service.SettingsService
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

	// init message broker
	messageBroker, err := rabbit.New(cfg)
	if err != nil {
		panic(err)
	}

	orchestrationService := orchestrator_service.New(
		log,
		operationSettingsStorage,
		operationSettingsStorage,
		messageBroker,
		cfg,
	)

	monitoringService := monitoring_service.New(
		log,
		cfg,
	)

	settingService := settings_service.New(
		log,
		operationSettingsStorage,
		cfg,
	)
	return &App{
		OrchestrationService: orchestrationService,
		MonitoringService:    monitoringService,
		SettingService:       settingService,
	}
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
