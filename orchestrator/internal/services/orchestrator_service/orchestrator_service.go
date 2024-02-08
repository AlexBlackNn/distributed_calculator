package orchestrator_service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"orchestrator/internal/config"
	"orchestrator/internal/domain/models"
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

	// if found that operation is in progress, that is result is nil
	if operationInDb.Operation != "" && operationInDb.Result == nil {
		return "", nil
	}

	execTimeModel, err := os.settingsStorage.GetOperationExecutionTime(ctx)
	execTime := message_broker.ExectutionTime{
		PlusOperationExecutionTime:           execTimeModel.PlusOperationExecutionTime,
		MinusOperationExecutionTime:          execTimeModel.MinusOperationExecutionTime,
		MultiplicationOperationExecutionTime: execTimeModel.MultiplicationExecutionTime,
		DivisionOperationExecutionTime:       execTimeModel.DivisionExecutionTime,
	}

	userId := uuid.New().String()
	message := message_broker.RequestMessage{
		Id:                    userId,
		MessageExectutionTime: execTime,
		Operation:             operation,
	}

	os.messageBroker.Send(ctx, message)

	operationModel := models.Operation{
		Id:        userId,
		Operation: operation,
	}

	os.operationStorage.SaveOperation(ctx, operationModel, nil)
	return userId, nil
}

func (os *OrchestratorService) ParseResponse(
	ctx context.Context,
) {
	// should return channel and using the channgel we need to write to postgres results
	fmt.Println("receiver")
	result, err := os.messageBroker.Receive()
	if err != nil {
		fmt.Println(err)
	}
	for msg := range result {
		fmt.Println(msg)

		opr := models.Operation{
			Id:     msg.Id,
			Result: msg.Value,
		}
		err := os.operationStorage.UpdateOperation(ctx, opr)
		if err != nil {
			fmt.Println(err)
		}
	}
}
