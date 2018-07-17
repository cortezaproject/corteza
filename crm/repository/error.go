package repository

type (
	repositoryError string
)

const (
	ErrDatabaseError = repositoryError("DatabaseError")
)

func (e repositoryError) Error() string {
	return "repository." + string(e)
}
