package store

import (
	"strings"

	"github.com/crusttech/human/server/pkg/errors"
)

type (
	// ErrorHandler
	// each implementation can have internal error handler that can translate
	// impl. specific errors like transaction
	ErrorHandler func(error) error
)

var (
	ErrNotFound  = errors.Plain(errors.KindNotFound, "not found")
	ErrNotUnique = errors.Plain(errors.KindDuplicateData, "not unique")
)

// ErrNotUniqueOn is ErrNotUnique with the offending resource and fields named.
//
// A bare "not unique" tells whoever hits it nothing: a resource can carry
// several unique constraints, and the service layer and the database can both
// raise it. Debugging then means guessing which one fired, and guessing wrong
// is expensive, because the two have completely different fixes.
//
// Same errors.KindDuplicateData as ErrNotUnique, so anything switching on the
// KIND keeps working; only the human-readable message gains detail.
// Returns *errors.Error, not error, so callers keep .Wrap()/.Stack() -- the
// driver error handlers wrap the underlying pq/mysql error onto it.
func ErrNotUniqueOn(resource string, fields ...string) *errors.Error {
	// errors.Plain formats its own message; pre-formatting with Sprintf would
	// pass a non-constant format string and mangle any % in a field name.
	if len(fields) == 0 {
		return errors.Plain(errors.KindDuplicateData, "%s not unique", resource)
	}

	return errors.Plain(
		errors.KindDuplicateData,
		"%s not unique: another one already uses this %s",
		resource,
		strings.Join(fields, " + "),
	)
}

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
