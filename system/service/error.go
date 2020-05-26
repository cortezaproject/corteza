package service

import (
	"github.com/pkg/errors"
)

type (
	serviceError string
)

const (
	ErrNoUpdatePermissions         serviceError = "NoUpdatePermissions"
	ErrNoReminderAssignPermissions serviceError = "NoReminderAssignPermissions"
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
