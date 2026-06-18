package dml

import (
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/system/types"
)

// fromDalModel converts a dal.Model to types.DmlModel.
func fromDalModel(m *dal.Model) *types.DmlModel {
	dm := &types.DmlModel{
		ConnectionID: m.ConnectionID,
		Ident:        m.Ident,
		Label:        m.Label,
		ResourceType: m.ResourceType,
	}
	for _, a := range m.Attributes {
		dm.Attributes = append(dm.Attributes, fromDalAttribute(a))
	}
	return dm
}

// fromDalAttribute converts a dal.Attribute to types.DmlAttribute.
func fromDalAttribute(a *dal.Attribute) *types.DmlAttribute {
	da := &types.DmlAttribute{
		Ident:      a.Ident,
		Label:      a.Label,
		PrimaryKey: a.PrimaryKey,
		Sortable:   a.Sortable,
		Filterable: a.Filterable,
	}

	if a.Type != nil {
		da.Type = string(a.Type.Type())
	}
	if a.Store != nil {
		da.Store = string(a.Store.Type())
	}

	return da
}
