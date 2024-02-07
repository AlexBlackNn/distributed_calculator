package storage

import (
	"context"
	"orchestrator/internal/domain/models"
	"time"
)

type OperationStorageInterface interface {
	SaveOperation(
		ctx context.Context,
		settings models.Operation,
	) error
	GetOperation(
		ctx context.Context,
		operation string,
	) (models.Operation, error)
}

type SettingsStorageInterface interface {
	SaveOperationExecutionTime(
		ctx context.Context,
		settings models.Settings,
	) error
	GetOperationExecutionTime(
		ctx context.Context,
	) (models.Settings, error)
}

type OperationCacheInterface interface {
	SaveOperation(ctx context.Context, operation string, result float64, ttl time.Duration) error
	GetOperation(ctx context.Context, operation string) (float64, error)
}
