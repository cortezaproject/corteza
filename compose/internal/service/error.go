package service

import (
	"github.com/pkg/errors"
)

type (
	serviceError string
)

const (
	ErrInvalidID           serviceError = "InvalidID"
	ErrStaleData           serviceError = "StaleData"
	ErrNoCreatePermissions serviceError = "NoCreatePermissions"
	ErrNoReadPermissions   serviceError = "NoReadPermissions"
	ErrNoUpdatePermissions serviceError = "NoUpdatePermissions"
	ErrNoDeletePermissions serviceError = "NoDeletePermissions"
	ErrNotImplemented      serviceError = "NotImplemented"
)

func (e serviceError) Error() string {
	return e.String()
}

func (e serviceError) String() string {
	return "crust.compose.service." + string(e)
}

func (e serviceError) withStack() error {
	return errors.WithStack(e)
}
