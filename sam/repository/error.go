package repository

type (
	repositoryError string
)

const (
	ErrDatabaseError  = repositoryError("DatabaseError")
	ErrNotImplemented = repositoryError("NotImplemented")
)

func (e repositoryError) Error() string {
	return e.String()
}

func (e repositoryError) String() string {
	return "crust.sam.repository." + string(e)
}
