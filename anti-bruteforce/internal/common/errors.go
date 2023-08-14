package common

import "errors"

var (
	ErrIPExists             = errors.New("the IP address already exists")
	ErrIPNotExists          = errors.New("the IP address is not exists")
	ErrUndefinedLoggerLevel = errors.New("the specified logger is not supported by the service")
)
