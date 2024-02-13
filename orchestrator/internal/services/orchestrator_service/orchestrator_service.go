package orchestrator_service

import (
	"context"
	"errors"
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
		fmt.Println("database Error", err)
	}
	// if found that operation is in progress (result is nil) returns saved id
	if operationInDb.Operation != "" {
		return operationInDb.Id, nil
	}

	execTimeModel, err := os.settingsStorage.GetOperationExecutionTime(ctx)
	execTime := message_broker.ExectutionTime{
		PlusOperationExecutionTime:           execTimeModel.PlusOperationExecutionTime,
		MinusOperationExecutionTime:          execTimeModel.MinusOperationExecutionTime,
		MultiplicationOperationExecutionTime: execTimeModel.MultiplicationExecutionTime,
		DivisionOperationExecutionTime:       execTimeModel.DivisionExecutionTime,
	}

	uid := uuid.New().String()
	message := message_broker.RequestMessage{
		Id:                    uid,
		MessageExectutionTime: execTime,
		Operation:             operation,
	}

	os.messageBroker.Send(ctx, message)

	operationModel := models.Operation{
		Id:        uid,
		Operation: operation,
	}

	os.operationStorage.SaveOperation(ctx, operationModel, nil)
	return uid, nil
}

func (os *OrchestratorService) ParseResponse(
	ctx context.Context,
) {
	// should return channel and using the channgel we need to write to postgres results
	fmt.Println("receiver")
	result, err := os.messageBroker.Receive()
	if err != nil {
		fmt.Println("+++++++++++++++", err)
	}
	for msg := range result {
		fmt.Println("====================>>>>", msg)
		var opr models.Operation
		if msg.Err != "" {
			opr = models.Operation{
				Id:     msg.Id,
				Result: msg.Value,
				Status: "Error",
			}
			fmt.Println("===========1111111=========>>>>", opr)
		} else {
			opr = models.Operation{
				Id:     msg.Id,
				Result: msg.Value,
				Status: "success",
			}
			fmt.Println("===========2222222=========>>>>", opr)
		}

		err := os.operationStorage.UpdateOperation(ctx, opr)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (os *OrchestratorService) CalculationResult(
	ctx context.Context,
	id string,
) (float64, error) {
	log := os.log.With(
		slog.String("info", "SERVICE LAYER: orchestrator_service.CalculationResult"),
	)
	print("111111111111111111111111")
	log.Info("check if operation was calculated")

	operationInDb, err := os.operationStorage.GetOperationById(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrOperationNotFound) {
			return 0, ErrNoOperation
		}
		return 0, fmt.Errorf("SERVICE LAYER: orchestrator_service.CalculationResult: %w", err)
	}
	// if found that operation is in progress (result is nil) returns saved id
	// TODO: NEED TO CHECK and FIX
	if nil == operationInDb.Result {
		// TODO: create as errors of service layer

		return 0, fmt.Errorf("Not Ready")
	}
	if operationInDb.Status == "Error" {
		// TODO: New error
		fmt.Println("=====))))))))))))))))))))))++++++++++++> ERROR")
		return 0, ErrFailedOperation
	}
	//TODO: create if it valid
	return operationInDb.Result.(float64), nil
}
