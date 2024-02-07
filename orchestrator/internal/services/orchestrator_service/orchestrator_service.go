package orchestrator_service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"orchestrator/internal/config"
	"orchestrator/message_broker"
	"orchestrator/storage"
)

type OrchestratorService struct {
	log *slog.Logger
	// data layer
	operationStorage storage.OperationStorageInterface
	settingsStorage  storage.SettingsStorageInterface
	// data layer
	operationCache storage.OperationCacheInterface
	messageBroker  message_broker.MessageBrokerInterface
	cfg            *config.Config
}

// New returns a new instance of Auth orchestrator_service
func New(
	log *slog.Logger,
	// data layer
	operationStorage storage.OperationStorageInterface,
	settingsStorage storage.SettingsStorageInterface,
	operationCache storage.OperationCacheInterface,
	// broker
	messageBroker message_broker.MessageBrokerInterface,

	cfg *config.Config,
) *OrchestratorService {
	return &OrchestratorService{
		log:              log,
		operationStorage: operationStorage,
		settingsStorage:  settingsStorage,
		operationCache:   operationCache,
		messageBroker:    messageBroker,
		cfg:              cfg,
	}
}

func (os *OrchestratorService) CalculationRequest(
	ctx context.Context,
	operation string,
) (string, error) {
	log := os.log.With(
		slog.String("info", "SERVICE LAYER: orchestrator_service.CalculationRequest"),
	)

	log.Info("check if operation was calculated")
	operationInCache, err := os.operationCache.GetOperation(ctx, operation)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(operationInCache)
	operationInDb, err := os.operationStorage.GetOperation(ctx, operation)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(operationInDb)

	//TODO: get settings from database
	execTime := message_broker.ExectutionTime{
		PlusOperationExecutionTime:           1000,
		MinusOperationExecutionTime:          2000,
		MultiplicationOperationExecutionTime: 1200,
		DivisionOperationExecutionTime:       2500,
	}

	userId := uuid.New().String()
	message := message_broker.RequestMessage{
		Id:                    userId,
		MessageExectutionTime: execTime,
		Operation:             operation,
	}

	os.messageBroker.Send(ctx, message)
	return userId, nil
}
