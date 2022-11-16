package store

import "github.com/cortezaproject/corteza/server/pkg/errors"

type (
	// each implementation can have internal error handler that can translate
	// impl. specific errors liker transaction
	ErrorHandler func(error) error
)

var (
	ErrNotFound  = errors.Plain(errors.KindNotFound, "not found")
	ErrNotUnique = errors.Plain(errors.KindDuplicateData, "not unique")
)

func HandleError(err error, h ErrorHandler) error {
	if err == nil {
		return nil
	}

	if h != nil {
		err = h(err)
	}

	if _, wrapped := err.(*errors.Error); wrapped {
		return err
	}

	return errors.
		Store("store error: %v", err).
		Apply(errors.StackSkip(1)).
		Wrap(err)
}
