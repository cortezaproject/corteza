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
	return "system.repository." + string(e)
}

func (e repositoryError) Eq(err error) bool {
	return err != nil && e.Error() == err.Error()
}
