package resource

import (
	"strconv"
)

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

	Identifiers []string
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
	ii := make(Identifiers, 0, len(ss))
	ii = ii.Add(ss...)
	return ii
}

func (ri Identifiers) Add(ii ...string) Identifiers {
	for _, i := range ii {
		if len(i) > 0 {
			ri = append(ri, i)
		}
	}

	return ri
}

func (ri Identifiers) Clone() Identifiers {
	out := make(Identifiers, 0, len(ri))
	for _, i := range ri {
		out = append(out, i)
	}

	return out
}

func (ri Identifiers) HasAny(check Identifiers) bool {
	// The size of these will be tiny so no need for hashmaps
	for _, i := range ri {
		for _, j := range check {
			if i == j {
				return true
			}
		}
	}

	return false
}

func (ri Identifiers) StringSlice() []string {
	return ri
}

func (ri Identifiers) First() string {
	if len(ri) == 0 {
		return ""
	}
	return ri[0]
}

func (ri Identifiers) FirstID() uint64 {
	if len(ri) <= 0 {
		return 0
	}

	for _, s := range ri {
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
	if len(r.Identifiers) == 0 {
		return false
	}

	for _, i := range r.Identifiers {
		if i == "*" {
			return true
		}
	}
	return false
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
			seen[r.ResourceType] = r.Identifiers.Clone()
			continue
		}

		// not yet seen
		if !ii.HasAny(r.Identifiers) {
			out = append(out, r)
			seen[r.ResourceType] = seen[r.ResourceType].Add(r.Identifiers...)
		}
	}

	return out
}
