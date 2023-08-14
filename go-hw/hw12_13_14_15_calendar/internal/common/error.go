package common

import "errors"

var (
	ErrUndefinedLoggerLevel = errors.New("the specified logger is not supported by the service")
	ErrNotFound             = errors.New("event not found")
	ErrDateBusy             = errors.New("this time is already occupied by another event")
)
