package service

import (
	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/types"
)

type (
	Mapper struct{}
)

// Merge copies the values from originating structure
// do the destination
//
// mostly, there will be less mapped fields on the destination
// side, so start looping from here
func (m *Mapper) Merge(in *ct.RecordValueSet, out *ct.RecordValueSet) {
	for _, destVal := range *out {
		for _, origVal := range *in {
			if destVal.Name == origVal.Name {
				destVal.Value = origVal.Value
				break
			}
		}
	}

	return
}

// Prepare creates a set of Records to be used later
// when the fields will be mapped via Merge()
func (m *Mapper) Prepare(mappings types.ModuleFieldMappingSet) (out ct.RecordValueSet) {
	for _, mm := range mappings {
		rv := &ct.RecordValue{
			Name:  mm.Destination.Name,
			Value: "",
		}
		out = append(out, rv)
	}

	return
}
