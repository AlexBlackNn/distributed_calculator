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
	messageBroker    message_broker.MessageBrokerInterface
	cfg              *config.Config
}

// New returns a new instance of orchestrator service
func New(
	log *slog.Logger,
	operationStorage storage.OperationStorageInterface,
	settingsStorage storage.SettingsStorageInterface,
	messageBroker message_broker.MessageBrokerInterface,
	cfg *config.Config,
) *OrchestratorService {
	return &OrchestratorService{
		log:              log,
		operationStorage: operationStorage,
		settingsStorage:  settingsStorage,
		messageBroker:    messageBroker,
		cfg:              cfg,
	}
}

// CalculationRequest starts calculation, sending data to message broker
func (os *OrchestratorService) CalculationRequest(
	ctx context.Context,
	operation string,
	appUser models.User,
) (string, error) {
	log := os.log.With(
		slog.String("info", "SERVICE LAYER: orchestrator_service.CalculationRequest"),
	)
	log.Info("check if operation was calculated")
	operationInDb, err := os.operationStorage.GetOperation(ctx, operation)

	if err != nil {
		if errors.Is(err, storage.ErrOperationNotFound) {
			log.Info("no operation in storage")

			log.Info("getting operation execution time from storage")
			var execTime models.Settings
			execTime, err = os.settingsStorage.GetOperationExecutionTime(ctx)
			if err != nil {
				log.Error("can't get operation execution time from storage")
				return "0", InternalError
			}
			log.Info("create message to worker")
			uid := uuid.New().String()
			message := message_broker.RequestMessage{
				Id:                    uid,
				MessageExectutionTime: execTime,
				Operation:             operation,
			}
			err = os.messageBroker.Send(ctx, message)
			if err != nil {
				log.Error("can't send data to message broker")
				return "0", InternalError
			}
			log.Info("save operation into storage")
			operationModel := models.Operation{
				Id:        uid,
				Operation: operation,
			}
			err = os.operationStorage.SaveOperation(ctx, operationModel, appUser, nil)
			if err != nil {
				log.Error("can't save operation to storage")
				return "0", InternalError
			}
			return uid, nil
		}
		return "0", InternalError
	}
	return operationInDb.Id, nil
}

// ParseResponse parses messages from message broker and write results to Storage
func (os *OrchestratorService) ParseResponse(
	ctx context.Context,
) {
	log := os.log.With(
		slog.String("info", "SERVICE LAYER: orchestrator_service.ParseResponse"),
	)
	log.Info("message broker receiver starts")

	result, err := os.messageBroker.Receive()
	if err != nil {
		log.Error("can't receive message ")
	}
	for msg := range result {
		log.Info("message received")
		var opr models.Operation
		if msg.Err != "" {
			opr = models.Operation{
				Id:     msg.Id,
				Result: msg.Value,
				Status: "Error",
			}
			log.Warn("status error in expression detected", "id", msg.Id, "result", msg.Value)
		} else {
			opr = models.Operation{
				Id:     msg.Id,
				Result: msg.Value,
				Status: "success",
			}
			log.Info("status success in expression detected", "id", msg.Id, "result", msg.Value)
		}

		err = os.operationStorage.UpdateOperation(ctx, opr)
		if err != nil {
			log.Error("can't write status operation into storage", "err", err, "id", msg.Id, "result", msg.Value)
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
	log.Info("check if operation was calculated")
	operationInDb, err := os.operationStorage.GetOperationById(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrOperationNotFound) {
			return 0, ErrNoOperation
		}
		return 0, fmt.Errorf("SERVICE LAYER: orchestrator_service.CalculationResult: %w", err)
	}

	if operationInDb.Status == "process" {
		log.Info("expression is in process")
		return 0, ErrOperationNotProcessed
	}
	if operationInDb.Status == "Error" {
		log.Warn("status error in expression detected")
		return 0, ErrFailedOperation
	}
	return operationInDb.Result.(float64), nil
}

func (os *OrchestratorService) GetOperationsWithPagination(
	ctx context.Context,
	pageSize int,
	cursor string,
) ([]models.Operation, error) {
	log := os.log.With(
		slog.String("info", "SERVICE LAYER: orchestrator_service.GetOperationsWithPagination"),
	)

	log.Info("Retrieving operations with pagination")

	operations, err := os.operationStorage.GetOperationsFastPagination(ctx, pageSize, cursor)
	if err != nil {
		return nil, fmt.Errorf("SERVICE LAYER: orchestrator_service.GetOperationsWithPagination: %w", err)
	}

	return operations, nil
}
