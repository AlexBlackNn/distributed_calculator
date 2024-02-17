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
	messageBroker, err := rabbit.New("test")
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

//func main() {
//	ctx := context.Background()
//
//	cfg := config.MustLoad()
//	// init logger
//	log := setupLogger(cfg.Env)
//	log.Info("starting application", slog.String("env", cfg.Env))
//
//	application := New(log, cfg)
//
//	fmt.Println("_______________________________________________________")
//	calculation := "http://guest:guest@localhost:15672/api/queues/%2f/operation"
//
//	// Send a GET request to the RabbitMQ Management API to get queue details
//	resp, err := grequests.Get(calculation, nil)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	var data map[string]interface{}
//	resp.JSON(&data)
//
//	if consumers, ok := data["consumers"].(float64); ok {
//		fmt.Println("Number of consumers connected to the queue:", consumers)
//	} else {
//		fmt.Println("Consumers information not available for this queue.")
//	}
//	fmt.Println("_______________________________________________________")
//	fmt.Println("*******************************************************")
//
//	id1, err := application.orchestrationService.CalculationRequest(ctx, "1*(1+(3*5))")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(">>>>>>>>>>>>> id1: ", id1)
//
//	go application.orchestrationService.ParseResponse(ctx)
//	id2, err := application.orchestrationService.CalculationRequest(ctx, "1*2+4/2")
//	fmt.Println(">>>>>>>>>>>>> id2: ", id1)
//
//	fmt.Println("*******************************************************")
//
//	time.Sleep(10 * time.Second)
//	fmt.Println("_______________________________________________________")
//	result1, err := application.orchestrationService.CalculationResult(ctx, id1)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(">>>>>>>>>>>>> result1: ", result1)
//	result2, err := application.orchestrationService.CalculationResult(ctx, id2)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(">>>>>>>>>>>>> result2: ", result2)
//	fmt.Println("_______________________________________________________")
//}

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
