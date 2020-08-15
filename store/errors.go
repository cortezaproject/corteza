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

var (
	ErrNotFound = &storeError{msg: "not found"}
)

func Error(msg string, err error) error {
	return &storeError{msg, err}
}
