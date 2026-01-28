package registry

import "errors"

var (
	ErrExecutableNotFound      = errors.New("registry: executable not found")
	ErrRevisionNotFound        = errors.New("registry: revision not found")
	ErrNoActiveRevisions       = errors.New("registry: no active revisions")
	ErrRevisionAlreadyExists   = errors.New("registry: revision already exists")
	ErrExecutableInUse         = errors.New("registry: executable is in use")
	ErrExecutableNotDeprecated = errors.New("registry: executable must be deprecated before removal")
	ErrUsageCheckerRequired    = errors.New("registry: usage checker required for removal")
)
