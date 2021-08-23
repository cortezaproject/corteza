package resource

import "strconv"

type (
	Interface interface {
		Identifiers() Identifiers
		ResourceType() string
		Refs() RefSet
		MarkPlaceholder()
		Placeholder() bool
		ReID(Identifiers)
		ReRef(old RefSet, new RefSet)
	}

	InterfaceSet []Interface

	IdentifiableInterface interface {
		Interface

		SysID() uint64
	}

	RefableInterface interface {
		Interface

		Ref() string
	}

	RBACInterface interface {
		Interface

		RBACParts() (string, *Ref, []*Ref)
	}

	LocaleInterface interface {
		Interface

		ResourceTranslationParts() (string, *Ref, []*Ref)
		EncodeTranslations() ([]*ResourceTranslation, error)
	}

	PrunableInterface interface {
		Interface

		Prune(*Ref)
	}

	RefSet []*Ref
	Ref    struct {
		// @todo check with Denis regarding strings here (the cdocs comment)
		// @todo should this become node type instead?
		ResourceType string
		Identifiers  Identifiers
		Constraints  RefSet
	}

	Identifiers map[string]bool
)

var (
	DataSourceResourceType  = "data:raw"
	SettingsResourceType    = "setting"
	RbacResourceType        = "rbac-rule"
	ResourceTranslationType = "resource-translation"
)

func MakeRef(rt string, ii Identifiers) *Ref {
	return &Ref{ResourceType: rt, Identifiers: ii}
}

func MakeWildRef(rt string) *Ref {
	return &Ref{ResourceType: rt, Identifiers: MakeIdentifiers("*")}
}

func MakeIdentifiers(ss ...string) Identifiers {
	ii := make(Identifiers)
	ii.Add(ss...)
	return ii
}

func (ri Identifiers) Add(ii ...string) Identifiers {
	for _, i := range ii {
		if len(i) > 0 {
			ri[i] = true
		}
	}

	return ri
}

func (ri Identifiers) HasAny(ii Identifiers) bool {
	for i := range ii {
		if ri[i] {
			return true
		}
	}

	return false
}

func (ri Identifiers) StringSlice() []string {
	ss := make([]string, 0, len(ri))
	for k := range ri {
		ss = append(ss, k)
	}
	return ss
}

func (ri Identifiers) First() string {
	ss := ri.StringSlice()
	if len(ss) <= 0 {
		return ""
	}
	return ss[0]
}

func (ri Identifiers) FirstID() uint64 {
	ss := ri.StringSlice()
	if len(ss) <= 0 {
		return 0
	}

	for _, s := range ss {
		if v, err := strconv.ParseUint(s, 10, 64); err != nil {
			continue
		} else {
			return v
		}
	}

	return 0
}

func (rr InterfaceSet) Walk(f func(r Interface) error) (err error) {
	for _, r := range rr {
		err = f(r)
		if err != nil {
			return
		}
	}

	return nil
}

// SearchForIdentifiers returns the resources where the provided identifiers exist
//
// The Resource is matching if at least one identifier matches.
func (rr InterfaceSet) SearchForIdentifiers(ii Identifiers) (out InterfaceSet) {
	out = make(InterfaceSet, 0, len(rr)/2)

	for _, r := range rr {
		if r.Identifiers().HasAny(ii) {
			out = append(out, r)
		}
	}

	return
}

// SearchForReferences returns the resources where the provided references exist
//
// The Resource is matching if at least one reference matches.
func (rr InterfaceSet) SearchForReferences(ref *Ref) (out InterfaceSet) {
	out = make(InterfaceSet, 0, len(rr)/2)

	for _, r := range rr {
		if r.Refs().HasRef(ref) {
			out = append(out, r)
		}
	}

	return
}

// Constraint returns the current reference with added constraint
func (r *Ref) Constraint(c *Ref) *Ref {
	if r.Constraints == nil {
		r.Constraints = make(RefSet, 0, 1)
	}

	r.Constraints = append(r.Constraints, &Ref{
		ResourceType: c.ResourceType,
		Identifiers:  MakeIdentifiers(c.Identifiers.StringSlice()...),
	})

	return r
}

// IsWildcard checks if this Ref points to all resources of a specific resource type
func (r *Ref) IsWildcard() bool {
	return r.Identifiers != nil && r.Identifiers["*"]
}

func (a *Ref) equals(b *Ref) bool {
	if a.ResourceType != b.ResourceType {
		return false
	}

	if !b.IsWildcard() && !a.Identifiers.HasAny(b.Identifiers) {
		return false
	}

	for _, c := range b.Constraints {
		if !a.Constraints.HasRef(c) {
			return false
		}
	}

	return true
}

func (rr RefSet) findRef(ref *Ref) int {
	for i, r := range rr {
		if r.equals(ref) {
			return i
		}
	}

	return -1
}

// replaceRef replaces the reference both on the ref level and on the
// constraint level.
func (rr RefSet) replaceRef(old, new *Ref) RefSet {
	found := false

	for x := len(rr) - 1; x >= 0; x-- {
		r := rr[x]

		if r.equals(old) {
			found = true
			if new == nil {
				rr[x] = rr[len(rr)-1]
				return rr[:len(rr)-1]
			} else {
				rr[x] = new
			}
		} else {
			if len(r.Constraints) > 0 {
				r.Constraints = r.Constraints.replaceRef(old, new)
			}
		}
	}

	if !found {
		return append(rr, new)
	}

	return rr
}

func (rr RefSet) HasRef(ref *Ref) bool {
	return rr.findRef(ref) > -1
}

// Unique returns only unique references
//
// Uniqueness is defined as "two references may not define
// the same resource type and identifier" combinations.
func (rr RefSet) Unique() RefSet {
	out := make(RefSet, 0, len(rr))
	seen := make(map[string]Identifiers)

	for _, r := range rr {
		ii, ok := seen[r.ResourceType]

		// type not seen at all, unique
		if !ok {
			out = append(out, r)
			seen[r.ResourceType] = r.Identifiers
			continue
		}

		// not yet seen
		if !ii.HasAny(r.Identifiers) {
			out = append(out, r)
			for i := range r.Identifiers {
				seen[r.ResourceType][i] = true
			}
		}
	}

	return out
}
