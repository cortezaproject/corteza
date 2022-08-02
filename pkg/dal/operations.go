package dal

type (
	Operation    string
	OperationSet []Operation
)

const (
	Create  Operation = "corteza::dal:operation:create"
	Update  Operation = "corteza::dal:operation:update"
	Delete  Operation = "corteza::dal:operation:delete"
	Search  Operation = "corteza::dal:operation:search"
	Lookup  Operation = "corteza::dal:operation:lookup"
	Paging  Operation = "corteza::dal:operation:paging"
	Sorting Operation = "corteza::dal:operation:sorting"
	Analyze Operation = "corteza::dal:operation:analyze"

	// @todo reporter operations
)

var (
	full = OperationSet{
		Create,
		Update,
		Delete,
		Search,
		Lookup,
		Paging,
		Sorting,
		Analyze,
	}

	createOperations = OperationSet{
		Create,
	}

	updateOperations = OperationSet{
		Update,
	}

	deleteOperations = OperationSet{
		Update,
	}

	searchOperations = OperationSet{
		Search,
		Paging,
		Sorting,
		Analyze,
	}

	lookupOperations = OperationSet{
		Lookup,
	}
)

// FullOperations returns all base system defined operations
func FullOperations() (cc OperationSet) {
	// Doing an union just to make a fresh copy
	return full.Union(nil)
}

// CreateOperations returns only requested operations used for Create operations
func CreateOperations(requested ...Operation) (required OperationSet) {
	return common(createOperations, requested)
}

// UpdateOperations returns only requested operations used for Update operations
func UpdateOperations(requested ...Operation) (required OperationSet) {
	return common(updateOperations, requested)
}

// DeleteOperations returns only requested operations used for delete operations
func DeleteOperations(requested ...Operation) (required OperationSet) {
	return common(deleteOperations, requested)
}

// SearchOperations returns only requested operations used for Search operations
func SearchOperations(requested ...Operation) (required OperationSet) {
	return common(searchOperations, requested)
}

// LookupOperations returns only requested operations used for Search operations
func LookupOperations(requested ...Operation) (required OperationSet) {
	return common(lookupOperations, requested)
}

func common(aa, bb OperationSet) OperationSet {
	return aa.Intersect(bb)
}

// ---

// IsSuperset is inverse IsSubset
//
// IsSuperset checks if all bb operations are inside aa
func (aa OperationSet) IsSuperset(bb ...Operation) bool {
	return OperationSet(bb).IsSubset(aa...)
}

// IsSubset checks if all aa operations are inside bb
func (aa OperationSet) IsSubset(bb ...Operation) bool {
	if len(aa) > len(bb) {
		return false
	}

	// When A is subset of B, the difference between the two must be 0
	return len(aa.Diff(bb)) == 0
}

// Intersect returns the intersection between the two sets
func (aa OperationSet) Intersect(bb OperationSet) (cc OperationSet) {
	cc = make(OperationSet, 0, len(aa))
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
func (aa OperationSet) Union(bb OperationSet) (cc OperationSet) {
	ix := make(map[Operation]bool)
	for _, c := range append(aa, bb...) {
		if !ix[c] {
			ix[c] = true
			cc = append(cc, c)
		}
	}
	return
}

// Diff calculates the difference between the two operation sets
//
// The diff uses aa as base
func (aa OperationSet) Diff(bb OperationSet) (cc OperationSet) {
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
