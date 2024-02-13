package settings_service

import (
	"context"
	"errors"
	"fmt"
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

func (ss *SettingsService) UpdateSettingsExecutionTime(
	ctx context.Context,
	operation_type string,
	execution_time int,
) error {
	var operatoion storage.OperationType

	switch operation_type {
	case "plus":
		operatoion = storage.PlusOperation
	case "minus":
		operatoion = storage.MinusOperation
	case "mult":
		operatoion = storage.MultiplicationOperation
	case "div":
		operatoion = storage.DivisionOperation
	default:
		//TODO: use storage errors
		return errors.New("Unknown operation type")
	}
	fmt.Print("dddddddddddddddddddddddddddddddd")
	log := ss.log.With(
		slog.String("info", "SERVICE LAYER: settings_service.PlusExecutionTime"),
	)
	log.Info("check if execution time is  valid")
	if !IsExecutionTimeValid(execution_time) {
		log.Info("execution time validation failed")
		return ErrValidationOperationTime
	}
	log.Info("execution time validation successful")
	fmt.Print("cccccccccccccccccccccccccccccccccccc")
	err := ss.settingsStorage.UpdateSettingsExecutionTime(ctx, operatoion, execution_time)
	if err != nil {
		log.Info("saving to storage execution time failed")
		return ErrServiceInternal
	}
	log.Info("execution time saved successfully")
	return nil
}

func IsExecutionTimeValid(execution_time int) bool {
	if execution_time > 0 && execution_time <= 10000 {
		return true
	}
	return false
}
