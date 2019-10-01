package service

import (
	"github.com/pkg/errors"
)

type (
	serviceError string
)

const (
	ErrInvalidID                      serviceError = "InvalidID"
	ErrInvalidHandle                  serviceError = "InvalidHandle"
	ErrNoPermissions                  serviceError = "NoPermissions"
	ErrNoGrantPermissions             serviceError = "NoGrantPermissions"
	ErrNoCreatePermissions            serviceError = "NoCreatePermissions"
	ErrNoUpdatePermissions            serviceError = "NoUpdatePermissions"
	ErrNoDeletePermissions            serviceError = "NoDeletePermissions"
	ErrNoReadPermissions              serviceError = "NoReadPermissions"
	ErrNoTriggerManagementPermissions serviceError = "NoTriggerManagementPermissions"
	ErrNoScriptCreatePermissions      serviceError = "NoScriptCreatePermissions"
	ErrNoReminderAssignPermissions    serviceError = "NoReminderAssignPermissions"
)

func (e serviceError) Error() string {
	return e.String()
}

func (e serviceError) String() string {
	return "system.service." + string(e)
}

func (e serviceError) withStack() error {
	return errors.WithStack(e)
}
