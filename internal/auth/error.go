package auth

import (
	"github.com/pkg/errors"
)

type (
	authError string
)

const (
	ErrConfigError = authError("ConfigError")
)

func (e authError) Error() string {
	return e.String()
}

func (e authError) String() string {
	return "internal.auth." + string(e)
}

func (e authError) New() error {
	return errors.WithStack(e)
}
