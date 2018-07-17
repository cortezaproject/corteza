package repository

type (
	repositoryError string
)

const (
	ErrDatabaseError  = repositoryError("DatabaseError")
	ErrNotImplemented = repositoryError("NotImplemented")
)

func (e repositoryError) Error() string {
	return "crust.sam.repository." + string(e)
}
