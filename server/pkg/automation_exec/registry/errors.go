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

func (r *registry) statsOnNewExecutable() {
	r.stats.TotalExecutables++
}

func (r *registry) statsOnRegisterRevision() {
	r.stats.TotalRevisions++
	r.stats.ActiveRevisions++
}

func (r *registry) statsOnDeprecateRevision() {
	r.stats.ActiveRevisions--
	r.stats.DeprecatedRevisions++
}

func (r *registry) statsOnRemoveDeprecatedRevision() {
	r.stats.TotalRevisions--
	r.stats.DeprecatedRevisions--
}

func (r *registry) statsOnRemoveExecutable() {
	r.stats.TotalExecutables--
}
