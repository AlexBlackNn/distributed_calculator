package orchestrator_service

import "errors"

var (
	ErrNoOperation           = errors.New("operation with requested uuid not found ")
	ErrFailedOperation       = errors.New("failed operation")
	ErrOperationNotProcessed = errors.New("operation not processed")
	InternalError            = errors.New("internal error")
)
