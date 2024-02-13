package storage

import "errors"

var (
	ErrUnknownOperationType = errors.New("unknown operation type")
	ErrOperationNotFound    = errors.New("operation not found")
)
