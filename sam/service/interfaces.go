package service

type (
	suspender interface {
		Suspend(ID uint64) error
		Unsuspend(ID uint64) error
	}

	archiver interface {
		Archive(ID uint64) error
		Unarchive(ID uint64) error
	}

	deleter interface {
		Delete(ID uint64) error
	}
)
