package dml

import (
	"context"
	"fmt"

	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/system/types"
)

type (
	// ComposeModuleSvc is the subset of compose module service that Applier needs.
	ComposeModuleSvc interface {
		FindByHandle(ctx context.Context, namespaceID uint64, handle string) (*composeTypes.Module, error)
		Create(ctx context.Context, mod *composeTypes.Module) (*composeTypes.Module, error)
		Update(ctx context.Context, mod *composeTypes.Module) (*composeTypes.Module, error)
	}

	// ComposeNamespaceSvc is the subset of compose namespace service that Applier and Importer need.
	ComposeNamespaceSvc interface {
		FindByHandle(ctx context.Context, handle string) (*composeTypes.Namespace, error)
		Create(ctx context.Context, ns *composeTypes.Namespace) (*composeTypes.Namespace, error)
	}

	// Applier creates or updates compose namespaces + modules from a DmlMapping.
	Applier struct {
		mapping   *Mapping
		moduleSvc ComposeModuleSvc
		nsSvc     ComposeNamespaceSvc
	}
)

func NewApplier(m *Mapping, mod ComposeModuleSvc, ns ComposeNamespaceSvc) *Applier {
	return &Applier{mapping: m, moduleSvc: mod, nsSvc: ns}
}

// Apply translates a persisted mapping into real compose namespaces + modules.
// Idempotent: existing modules (matched by handle) are updated in-place.
func (a *Applier) Apply(ctx context.Context, mappingID uint64) error {
	mp, err := a.mapping.FindByID(ctx, mappingID)
	if err != nil {
		return err
	}

	ns, err := a.ensureNamespace(ctx, mp)
	if err != nil {
		return fmt.Errorf("dml apply: namespace: %w", err)
	}

	for _, tbl := range mp.Tables {
		if tbl.Skip {
			continue
		}
		if err := a.applyTable(ctx, ns, tbl, mp); err != nil {
			return fmt.Errorf("dml apply: table %q: %w", tbl.SourceIdent, err)
		}
	}

	return nil
}

func (a *Applier) ensureNamespace(ctx context.Context, mp *types.DmlMapping) (*composeTypes.Namespace, error) {
	handle := mp.NamespaceHandle
	if handle == "" {
		handle = fmt.Sprintf("dml_%d", mp.ConnectionID)
	}

	ns, err := a.nsSvc.FindByHandle(ctx, handle)
	if err == nil {
		return ns, nil
	}

	return a.nsSvc.Create(ctx, &composeTypes.Namespace{
		Slug:    handle,
		Name:    handle,
		Enabled: true,
	})
}

func (a *Applier) applyTable(ctx context.Context, ns *composeTypes.Namespace, tbl *types.DmlTableMap, mp *types.DmlMapping) error {
	fields := make(composeTypes.ModuleFieldSet, 0, len(tbl.Columns))
	for _, col := range tbl.Columns {
		if col.Skip {
			continue
		}
		fields = append(fields, &composeTypes.ModuleField{
			NamespaceID: ns.ID,
			Name:        col.FieldName,
			Label:       col.FieldName,
			Kind:        col.FieldKind,
		})
	}

	existing, err := a.moduleSvc.FindByHandle(ctx, ns.ID, tbl.ModuleHandle)
	if err == nil {
		// update
		existing.Fields = fields
		existing.Name = tbl.ModuleName
		if existing.Name == "" {
			existing.Name = tbl.ModuleHandle
		}
		_, err = a.moduleSvc.Update(ctx, existing)
		return err
	}

	name := tbl.ModuleName
	if name == "" {
		name = tbl.ModuleHandle
	}

	_, err = a.moduleSvc.Create(ctx, &composeTypes.Module{
		Handle:      tbl.ModuleHandle,
		Name:        name,
		NamespaceID: ns.ID,
		Fields:      fields,
	})
	return err
}
