package tmp

import "github.com/cortezaproject/corteza-server/pkg/envoy/resource"

type (
	importState struct {
		state resNodeMap
		// We'll use this to keep track of the existing resources
		existing map[resource.Interface]uint64
	}

	// ref: internal ID
	refMap map[string]uint64
	// resource type: refMap
	resMap     map[string]refMap
	resNodeMap map[resource.Interface]resMap
)

func NewImportState() *importState {
	return &importState{
		state:    make(resNodeMap),
		existing: make(map[resource.Interface]uint64),
	}
}

func (s *importState) AddRefMapping(r resource.Interface, res string, id uint64, refs ...string) {
	if s.state[r] == nil {
		s.state[r] = make(resMap)
	}

	if s.state[r][res] == nil {
		s.state[r][res] = make(refMap)
	}

	for _, ref := range refs {
		s.state[r][res][ref] = id
	}
}

func (s *importState) Existint(r resource.Interface) uint64 {
	return s.existing[r]
}

func (s *importState) AddExisting(r resource.Interface, resID uint64) {
	s.existing[r] = resID
}
