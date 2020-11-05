package tmp

import "github.com/cortezaproject/corteza-server/pkg/envoy/resource"

type (
	encoderState struct {
		// Maps the state for each separate resource
		state map[resource.Interface]resRefs
		// Lets us keep track of existing resources
		// @todo how will we deal with CompseRecord?
		existing map[resource.Interface]bool
	}

	// ref: internalID
	ref map[string]uint64
	// resourceType: ref
	resRefs map[string]ref
)

// NewEncoderState initializes and returns an empty encoder state
func NewEncoderState() *encoderState {
	return &encoderState{
		state:    make(map[resource.Interface]resRefs),
		existing: make(map[resource.Interface]bool),
	}
}

// Get returns the state for the given resource
func (es *encoderState) Get(res resource.Interface) resRefs {
	return es.state[res]
}

// Set sets the encoding state for the given resource
func (es *encoderState) Set(res resource.Interface, resType string, id uint64, refs ...string) {
	if es.state[res] == nil {
		es.state[res] = make(resRefs)
	}

	es.state[res].Set(resType, id, refs...)
}

// Merge merges the encoding state with another state
func (es *encoderState) Merge(res resource.Interface, refs resRefs) {
	if refs == nil {
		return
	}

	if es.state[res] == nil {
		es.state[res] = make(resRefs)
	}

	es.state[res].Merge(refs)
}

// Set sets refs for a given resource
func (r resRefs) Set(res string, id uint64, refs ...string) {
	if r[res] == nil {
		r[res] = make(ref)
	}

	r[res].Set(id, refs...)
}

// Merge merges the two states
func (r resRefs) Merge(state resRefs) {
	for res, refs := range state {
		if r[res] == nil {
			r[res] = make(ref)
		}

		for lID, iID := range refs {
			r[res][lID] = iID
		}
	}
}

// Get returns an internalID based on the passed interface.
//
// If not found, it returns a 0
func (r resRefs) Get(res resource.Interface) uint64 {
	ref := r[res.ResourceType()]
	if ref == nil {
		return 0
	}

	for i := range res.Identifiers() {
		if ref[i] > 0 {
			return ref[i]
		}
	}
	return 0
}

// Set sets provided ref: ID mappings
func (r ref) Set(id uint64, refs ...string) {
	for _, ref := range refs {
		r[ref] = id
	}
}

func (es *encoderState) Exists(res resource.Interface) bool {
	return es.existing[res]
}

func (es *encoderState) SetExists(res resource.Interface) {
	es.existing[res] = true
}
