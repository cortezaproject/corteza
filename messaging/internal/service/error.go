package service

type (
	readableError string
)

func (e readableError) Error() string {
	return string(e)
}

const (
	ErrNoPermissions readableError = "You don't have permissions for this operation"
)
