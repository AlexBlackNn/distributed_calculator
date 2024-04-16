package utils

import "errors"

var (
	ErrNoJWT = errors.New("jwt token absent")
)
