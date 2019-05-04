package repository

import (
	"github.com/pkg/errors"
)

type (
	repositoryError string
)

const (
	ErrNotImplemented = repositoryError("NotImplemented")
)

func (e repositoryError) Error() string {
	return e.String()
}

func (e repositoryError) String() string {
	return "crust.compose.repository." + string(e)
}

func (e repositoryError) new() error {
	return errors.WithStack(e)
}
