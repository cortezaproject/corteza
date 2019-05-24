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
	ErrNoGrantPermissions  serviceError = "NoGrantPermissions"
	ErrNoCreatePermissions serviceError = "NoCreatePermissions"
	ErrNoReadPermissions   serviceError = "NoReadPermissions"
	ErrNoUpdatePermissions serviceError = "NoUpdatePermissions"
	ErrNoDeletePermissions serviceError = "NoDeletePermissions"
	ErrNamespaceRequired   serviceError = "NamespaceRequired"
	ErrModulePageExists    serviceError = "ModulePageExists"
	ErrNotImplemented      serviceError = "NotImplemented"
)

func (e serviceError) Error() string {
	return e.String()
}

func (e serviceError) String() string {
	return "compose.service." + string(e)
}

func (e serviceError) withStack() error {
	return errors.WithStack(e)
}
