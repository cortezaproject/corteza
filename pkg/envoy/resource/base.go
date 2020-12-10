package resource

type (
	base struct {
		rt string
		ii Identifiers
		rr RefSet

		ts    *Timestamps
		us    *Userstamps
		cfg   *EnvoyConfig
		urefs RefSet
	}

	EnvoyConfig struct {
		// SkipIf determines when the encoding should be skipped for this resource
		SkipIf     string
		OnExisting MergeAlg
	}

	Timestamps struct {
		CreatedAt   string
		UpdatedAt   string
		DeletedAt   string
		ArchivedAt  string
		SuspendedAt string
	}
	Userstamps struct {
		CreatedBy string
		UpdatedBy string
		DeletedBy string
		OwnedBy   string
	}

	MergeAlg int
)

const (
	// Default takes the operation defined default
	Default MergeAlg = iota
	// Skip skips the existing resource
	Skip
	// Replace replaces the existing resource
	Replace
	// MergeLeft updates the existing resource, giving priority to the existing data
	MergeLeft
	// MergeRight updates the existing resource, giving priority to the new data
	MergeRight
)

// State management methods

// AddIdentifier adds a set of identifiers to the current resource
func (t *base) AddIdentifier(ss ...string) {
	if t.ii == nil {
		t.ii = make(Identifiers)
	}

	t.ii.Add(ss...)
}

// AddRef adds a new reference to the current resource
func (t *base) AddRef(rt string, ii ...string) *Ref {
	if t.rr == nil {
		t.rr = make(RefSet, 0, 10)
	}

	iiC := make([]string, 0, len(ii))
	for _, i := range ii {
		if i != "" {
			iiC = append(iiC, i)
		}
	}

	ref := &Ref{ResourceType: rt, Identifiers: Identifiers{}.Add(iiC...)}
	t.rr = append(t.rr, ref)

	return ref
}

// SetResourceType sets the resource type of the current resource struct
func (t *base) SetResourceType(rt string) {
	t.rt = rt
}

func (t *base) SetTimestamps(ts *Timestamps) {
	t.ts = ts
}
func (t *base) Timestamps() *Timestamps {
	return t.ts
}

func (t *base) SetUserstamps(us *Userstamps) {
	t.us = us

	if us != nil {
		uu := []string{us.CreatedBy, us.UpdatedBy, us.DeletedBy, us.OwnedBy}
		t.SetUserRefs(uu)
	}
}
func (t *base) Userstamps() *Userstamps {
	return t.us
}

func (t *base) SetConfig(cfg *EnvoyConfig) {
	t.cfg = cfg
}
func (t *base) Config() *EnvoyConfig {
	return t.cfg
}

func (t *base) SetUserRefs(uu []string) {
	if t.urefs == nil {
		t.urefs = make(RefSet, 0, 4)
	}

	for _, u := range uu {
		if u != "" {
			t.urefs = append(t.urefs, t.AddRef(USER_RESOURCE_TYPE, u))
		}
	}
}
func (t *base) UserRefs() RefSet {
	return t.urefs
}

// Resource interface methods

func (t *base) Identifiers() Identifiers {
	return t.ii
}
func (t *base) ResourceType() string {
	return t.rt
}
func (t *base) Refs() RefSet {
	return t.rr
}
