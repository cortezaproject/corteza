package resource

type (
	base struct {
		rt string
		ii Identifiers
		rr RefSet
	}
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
func (t *base) AddRef(rt string, ii ...string) {
	if t.rr == nil {
		t.rr = make(RefSet, 0, 10)
	}

	ref := &Ref{ResourceType: rt, Identifiers: Identifiers{}.Add(ii...)}
	t.rr = append(t.rr, ref)
}

// SetResourceType sets the resource type of the current resource struct
func (t *base) SetResourceType(rt string) {
	t.rt = rt
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
