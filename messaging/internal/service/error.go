package service

import (
	"github.com/pkg/errors"
)

type (
	readableError string
)

func (e readableError) Error() string {
	return string(e)
}

func (e readableError) new() error {
	return errors.WithStack(e)
}

const (
	ErrNoPermissions  readableError = "You don't have permissions for this operation"
	ErrUnknownCommand readableError = "Unknown command"
)
