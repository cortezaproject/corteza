package proto

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

func FromNamespace(i *types.Namespace) *Namespace {
	if i == nil {
		return nil
	}

	var p = &Namespace{
		NamespaceID: i.ID,
		Name:        i.Name,
		Slug:        i.Slug,
		Enabled:     i.Enabled,
		CreatedAt:   fromTime(i.CreatedAt),
		UpdatedAt:   fromTime(i.UpdatedAt),
		DeletedAt:   fromTime(i.DeletedAt),
	}

	return p
}
