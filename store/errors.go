package store

type (
	storeError struct {
		msg string
		err error
	}
)

func (e *storeError) Error() string {
	return e.msg
}

func (e *storeError) Unwrap() error {
	return e.err
}

func (e *storeError) Wrap(err error) error {
	e.err = err
	return e
}

var (
	ErrNotFound  = &storeError{msg: "not found"}
	ErrNotUnique = &storeError{msg: "not unique"}
)

func WrappedError(err error, msg string, a ...interface{}) error {
	return &storeError{msg, err}
}
