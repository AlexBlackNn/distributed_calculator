package orchestrator_service

import "errors"

var (
	ErrNoOperation     = errors.New("operation with requested uuid not found ")
	ErrFailedOperation = errors.New("failed operation")
)
