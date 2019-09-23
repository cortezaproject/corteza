package proto

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

func FromModule(i *types.Module) *Module {
	if i == nil {
		return nil
	}

	var p = &Module{
		ModuleID:    i.ID,
		NamespaceID: i.NamespaceID,
		Name:        i.Name,
		CreatedAt:   fromTime(i.CreatedAt),
		UpdatedAt:   fromTime(i.UpdatedAt),
		DeletedAt:   fromTime(i.DeletedAt),
		Fields:      make([]*ModuleField, len(i.Fields)),
	}

	for f := range i.Fields {
		p.Fields[f] = &ModuleField{
			FieldID: i.Fields[f].ID,
			Name:    i.Fields[f].Name,
			Kind:    i.Fields[f].Kind,

			IsMulti:    i.Fields[f].Multi,
			IsPrivate:  i.Fields[f].Private,
			IsVisible:  i.Fields[f].Visible,
			IsRequired: i.Fields[f].Required,

			// @todo add default value
			// @todo add options
		}
	}

	return p
}
