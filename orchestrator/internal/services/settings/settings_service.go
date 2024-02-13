package settings_service

import (
	"context"
	"log/slog"
	"orchestrator/internal/config"
	"orchestrator/storage"
)

type SettingsService struct {
	log             *slog.Logger
	settingsStorage storage.SettingsStorageInterface
	cfg             *config.Config
}

// New returns a new instance of Settings Service
func New(
	log *slog.Logger,
	// data layer
	settingsStorage storage.SettingsStorageInterface,
	cfg *config.Config,
) *SettingsService {
	return &SettingsService{
		log:             log,
		settingsStorage: settingsStorage,
		cfg:             cfg,
	}
}

func (os *SettingsService) PlusExecutionTime(
	ctx context.Context,
	execution_time int,
) error {
	//log := os.log.With(
	//	slog.String("info", "SERVICE LAYER: settings_service.PlusExecutionTime"),
	//)

	//log.Info("check if operation was calculated")
	//operationInCache, err := os.operationCache.GetOperation(ctx, operation)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(operationInCache)
	//operationInDb, err := os.operationStorage.GetOperation(ctx, operation)
	//if err != nil {
	//	fmt.Println("database Error", err)
	//}
	//// if found that operation is in progress (result is nil) returns saved id
	//if operationInDb.Operation != "" {
	//	return operationInDb.Id, nil
	//}
	//
	//execTimeModel, err := os.settingsStorage.GetOperationExecutionTime(ctx)
	//execTime := message_broker.ExectutionTime{
	//	PlusOperationExecutionTime:           execTimeModel.PlusOperationExecutionTime,
	//	MinusOperationExecutionTime:          execTimeModel.MinusOperationExecutionTime,
	//	MultiplicationOperationExecutionTime: execTimeModel.MultiplicationExecutionTime,
	//	DivisionOperationExecutionTime:       execTimeModel.DivisionExecutionTime,
	//}
	//
	//uid := uuid.New().String()
	//message := message_broker.RequestMessage{
	//	Id:                    uid,
	//	MessageExectutionTime: execTime,
	//	Operation:             operation,
	//}
	//
	//os.messageBroker.Send(ctx, message)
	//
	//operationModel := models.Operation{
	//	Id:        uid,
	//	Operation: operation,
	//}

	//os.operationStorage.SaveOperation(ctx, operationModel, nil)
	return nil
}
