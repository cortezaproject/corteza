package registry

type (
	stats struct {
		TotalExecutables    int
		TotalRevisions      int
		ActiveRevisions     int
		DeprecatedRevisions int
	}
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
