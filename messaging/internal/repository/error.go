package repository

import (
	"github.com/pkg/errors"
)

type (
	repositoryError string
)

const (
	ErrDatabaseError    = repositoryError("DatabaseError")
	ErrNotImplemented   = repositoryError("NotImplemented")
	ErrConfigError      = repositoryError("ConfigError")
	ErrEventsPullClosed = repositoryError("EventsPullClosed")
)

func (e repositoryError) Error() string {
	return e.String()
}

func (e repositoryError) String() string {
	return "messaging.repository." + string(e)
}

func (e repositoryError) New() error {
	return errors.WithStack(e)
}
