package types

type (
	Flag struct {
		// Kind of the flagged resource
		Kind string
		// ID if the flagged resource
		ResourceID uint64
		// The owner of this flag; 0 = everyone
		OwnedBy uint64

		Name   string
		Active bool
	}

	// @todo codegen this thing
	FlagSet []*Flag

	FlagFilter struct {
		Kind       string
		ResourceID []uint64
		OwnedBy    []uint64
		Name       []string
		Limit      uint
	}
)

const (
	FlagResourceType = "corteza::generic:flag"
)
