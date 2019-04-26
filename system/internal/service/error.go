package service

import (
	"github.com/pkg/errors"
)

type (
	serviceError  string
	readableError string
)

func (e serviceError) Error() string {
	return "crust.messaging.service." + string(e)
}

func (e readableError) Error() string {
	return string(e)
}

func (e readableError) new() error {
	return errors.WithStack(e)
}

const (
	ErrNoPermissions   readableError = "You don't have permissions for this operation"
	ErrAvatarOnlyHTTPS readableError = "Avatar URL only supports HTTPS"
)
