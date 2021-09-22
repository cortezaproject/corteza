package resource

import (
	"fmt"
	"strconv"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/system/types"
)

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

// fn converts identifier values (string, fmt.Stringer, uint64) to string slice
//
// Each value is checked and should not be empty or zero
func identifiers(ii ...interface{}) []string {
	ss := make([]string, 0, len(ii))

	for _, i := range ii {
		switch c := i.(type) {
		case uint64:
			if c == 0 {
				continue
			}

			ss = append(ss, strconv.FormatUint(c, 10))

		case fmt.Stringer:
			if c.String() == "" {
				continue
			}

			ss = append(ss, c.String())

		case string:
			if c == "" {
				continue
			}

			ss = append(ss, c)
		}
	}

	return ss
}

// Check checks if the identifier c is in the set of identifiers ii
func Check(a string, ii ...interface{}) bool {
	for _, i := range ii {
		switch pi := i.(type) {
		case string:
			if a == pi {
				return true
			}
		case uint64:
			if pi > 0 && strconv.FormatUint(pi, 10) == i {
				return true
			}
		}
	}

	return false
}

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
		uu := []*Userstamp{us.CreatedBy, us.UpdatedBy, us.DeletedBy, us.OwnedBy, us.RunAs}
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

func (t *base) SetUserRefs(uu []*Userstamp) {
	if t.urefs == nil {
		t.urefs = make(RefSet, 0, 4)
	}

	for _, u := range uu {
		if u == nil {
			continue
		}
		if u.UserID > 0 {
			t.urefs = append(t.urefs, t.AddRef(types.UserResourceType, strconv.FormatUint(u.UserID, 10)))
		} else if u.Ref != "" {
			t.urefs = append(t.urefs, t.AddRef(types.UserResourceType, u.Ref))
		}
	}
}

func (t *base) UserRefs() RefSet {
	return t.urefs
}

func (t *base) Identifiers() Identifiers {
	return t.ii
}

func (t *base) ResourceType() string {
	return t.rt
}

func (t *base) Refs() RefSet {
	return t.rr
}

func (t *base) Ref() *Ref {
	return &Ref{ResourceType: t.rt, Identifiers: t.ii}
}

func (t *base) HasRefs() bool {
	return t.rr == nil || len(t.rr) == 0
}

func IgnoreDepResolution(ref *Ref) bool {
	return ref.ResourceType == composeTypes.ModuleFieldResourceType
}
