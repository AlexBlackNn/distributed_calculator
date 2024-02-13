package settings_service

import "errors"

var (
	ErrServiceInternal         = errors.New("internal error ")
	ErrValidationOperationTime = errors.New("failed operation")
)
