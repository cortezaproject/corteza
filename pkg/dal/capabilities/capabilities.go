package capabilities

type (
	Capability string
	Set        []Capability
)

const (
	Create  Capability = "corteza::dal:capability:create"
	Update  Capability = "corteza::dal:capability:update"
	Delete  Capability = "corteza::dal:capability:delete"
	Search  Capability = "corteza::dal:capability:search"
	Lookup  Capability = "corteza::dal:capability:lookup"
	Paging  Capability = "corteza::dal:capability:paging"
	Stats   Capability = "corteza::dal:capability:stats"
	Sorting Capability = "corteza::dal:capability:sorting"
	RBAC    Capability = "corteza::dal:capability:RBAC"
)

var (
	full = Set{
		Create,
		Update,
		Search,
		Lookup,
		Paging,
		Stats,
		Sorting,
		RBAC,
	}

	accessControlCapabilities = Set{
		RBAC,
	}

	createCapabilities = Set{
		RBAC,
		Create,
	}

	updateCapabilities = Set{
		RBAC,
		Update,
	}

	deleteCapabilities = Set{
		RBAC,
		Update,
	}

	searchCapabilities = Set{
		Search,
		Paging,
		Sorting,
		Stats,
		RBAC,
	}

	lookupCapabilities = Set{
		Lookup,
		RBAC,
	}
)

// FullCapabilities returns all base system defined capabilities
func FullCapabilities() (cc Set) {
	// Doing an union just to make a fresh copy
	return full.Union(nil)
}

// AccessControlCapabilities returns only requested capabilities used for AccessControl operations
func AccessControlCapabilities(requested ...Capability) (required Set) {
	return common(accessControlCapabilities, requested)
}

// CreateCapabilities returns only requested capabilities used for Create operations
func CreateCapabilities(requested ...Capability) (required Set) {
	return common(createCapabilities, requested)
}

// UpdateCapabilities returns only requested capabilities used for Update operations
func UpdateCapabilities(requested ...Capability) (required Set) {
	return common(updateCapabilities, requested)
}

// DeleteCapabilities returns only requested capabilities used for delete operations
func DeleteCapabilities(requested ...Capability) (required Set) {
	return common(deleteCapabilities, requested)
}

// SearchCapabilities returns only requested capabilities used for Search operations
func SearchCapabilities(requested ...Capability) (required Set) {
	return common(searchCapabilities, requested)
}

// LookupCapabilities returns only requested capabilities used for Search operations
func LookupCapabilities(requested ...Capability) (required Set) {
	return common(lookupCapabilities, requested)
}

func common(aa, bb Set) Set {
	return aa.Intersect(bb)
}

// ---

// IsSuperset is inverse IsSubset
//
// IsSuperset checks if all bb capabilities are inside aa
func (aa Set) IsSuperset(bb ...Capability) bool {
	return Set(bb).IsSubset(aa...)
}

// IsSubset checks if all aa capabilities are inside bb
func (aa Set) IsSubset(bb ...Capability) bool {
	if len(aa) > len(bb) {
		return false
	}

	// When A is subset of B, the difference between the two must be 0
	return len(aa.Diff(bb)) == 0
}

// Intersect returns the intersection between the two sets
func (aa Set) Intersect(bb Set) (cc Set) {
	cc = make(Set, 0, len(aa))
	for _, a := range aa {
		for _, b := range bb {
			if a == b {
				cc = append(cc, a)
				break
			}
		}
	}

	return
}

// Intersect returns the union between the two sets
//
// Duplicates are omitted
func (aa Set) Union(bb Set) (cc Set) {
	ix := make(map[Capability]bool)
	for _, c := range append(aa, bb...) {
		if !ix[c] {
			ix[c] = true
			cc = append(cc, c)
		}
	}
	return
}

// Diff calculates the difference between the two capability sets
//
// The diff uses aa as base
func (aa Set) Diff(bb Set) (cc Set) {
	for _, a := range aa {
		found := false
		for _, b := range bb {
			found = a == b
			if found {
				break
			}
		}

		if found {
			continue
		}
		cc = append(cc, a)
	}

	return
}
