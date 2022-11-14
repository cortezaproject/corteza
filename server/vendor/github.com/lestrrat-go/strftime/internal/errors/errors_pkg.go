//go:build !strftime_native_errors
// +build !strftime_native_errors

package errors

import "github.com/pkg/errors"

func New(s string) error {
	return errors.New(s)
}

func Errorf(s string, args ...interface{}) error {
	return errors.Errorf(s, args...)
}

func Wrap(err error, s string) error {
	return errors.Wrap(err, s)
}
