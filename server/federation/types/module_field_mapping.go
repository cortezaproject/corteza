package types

const (
	ModuleFieldMappingSetFindTypeOrigin ModuleFieldMappingSetFindType = iota
	ModuleFieldMappingSetFindTypeDestination
)

type (
	ModuleFieldMappingSet []*ModuleFieldMapping

	ModuleFieldMappingSetFindType int
)

// Find looks up a origin or destination mapping
func (list *ModuleFieldMappingSet) FindByName(name string, findType ModuleFieldMappingSetFindType) (*ModuleFieldMapping, error) {
	for _, mfm := range *list {
		switch findType {

		case ModuleFieldMappingSetFindTypeOrigin:
			if mfm.Origin.Name == name {
				return mfm, nil
			}

		case ModuleFieldMappingSetFindTypeDestination:
			if mfm.Destination.Name == name {
				return mfm, nil
			}
		}
	}

	return nil, nil
}
