package service

import (
	"github.com/pkg/errors"
)

type (
	serviceError string
)

const (
	ErrInvalidID     serviceError = "InvalidID"
	ErrNoPermissions serviceError = "NoPermissions"
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
