package repository

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
	return "compose.repository." + string(e)
}
