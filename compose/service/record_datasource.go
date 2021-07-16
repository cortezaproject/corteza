package service

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/report"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/spf13/cast"
)

func (svc record) Datasource(ctx context.Context, ld *report.LoadStepDefinition) (report.Datasource, error) {
	var (
		moduleID    uint64
		namespaceID uint64
		err         error

		def = ld.Definition
	)

	if mr, ok := def["namespaceID"]; ok {
		namespaceID, err = cast.ToUint64E(mr)
		if err != nil {
			return nil, err
		}
	} else if mr, ok = def["namespace"]; ok {
		// slug; fetch from store
		ns, err := store.LookupComposeNamespaceBySlug(ctx, svc.store, mr.(string))
		if errors.IsNotFound(err) {
			err = NamespaceErrNotFound()
		}
		if err != nil {
			return nil, err
		}
		namespaceID = ns.ID
	} else {
		return nil, fmt.Errorf("compose namespace not defined")
	}

	if mr, ok := def["moduleID"]; ok {
		moduleID, err = cast.ToUint64E(mr)
		if err != nil {
			return nil, err
		}
	} else if mr, ok = def["module"]; ok {
		// handle; fetch from store
		mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, svc.store, namespaceID, mr.(string))
		if errors.IsNotFound(err) {
			err = ModuleErrNotFound()
		}
		if err != nil {
			return nil, err
		}
		moduleID = mod.ID
	}

	// Find mod
	mod, err := svc.store.LookupComposeModuleByID(ctx, moduleID)
	if err != nil {
		return nil, err
	}
	mod.Fields, _, err = svc.store.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{
		ModuleID: []uint64{mod.ID},
	})
	if err != nil {
		return nil, err
	}

	if len(ld.Columns) == 0 {
		cols := make(report.FrameColumnSet, 0, len(mod.Fields)+8)

		var c *report.FrameColumn
		c = report.MakeColumnOfKind("Record")
		c.Name = "id"
		c.Label = "Record ID"
		cols = append(cols, c)

		for _, f := range mod.Fields {
			k := f.Kind
			c = report.MakeColumnOfKind(k)
			c.Name = f.Name
			c.Label = f.Label
			if c.Label == "" {
				c.Label = c.Name
			}
			cols = append(cols, c)
		}

		// Sys fields
		c = report.MakeColumnOfKind("Date")
		c.Name = "createdAt"
		c.Label = "Created at"
		cols = append(cols, c)

		c = report.MakeColumnOfKind("User")
		c.Name = "createdBy"
		c.Label = "Created by"
		cols = append(cols, c)

		c = report.MakeColumnOfKind("Date")
		c.Name = "updatedAt"
		c.Label = "Updated at"
		cols = append(cols, c)

		c = report.MakeColumnOfKind("User")
		c.Name = "updatedBy"
		c.Label = "Updated by"
		cols = append(cols, c)

		c = report.MakeColumnOfKind("Date")
		c.Name = "deletedAt"
		c.Label = "Deleted at"
		cols = append(cols, c)

		c = report.MakeColumnOfKind("User")
		c.Name = "deletedBy"
		c.Label = "Deleted by"
		cols = append(cols, c)

		ld.Columns = cols
	}

	return svc.store.ComposeRecordDatasource(ctx, mod, ld)
}
