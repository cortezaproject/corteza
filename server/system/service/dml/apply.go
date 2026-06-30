package dml

import (
	"context"
	"fmt"

	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/system/types"
)

type (
	ComposeModuleSvc interface {
		FindByHandle(ctx context.Context, namespaceID uint64, handle string) (*composeTypes.Module, error)
		Create(ctx context.Context, mod *composeTypes.Module) (*composeTypes.Module, error)
		Update(ctx context.Context, mod *composeTypes.Module) (*composeTypes.Module, error)
	}

	ComposeNamespaceSvc interface {
		FindByHandle(ctx context.Context, handle string) (*composeTypes.Namespace, error)
		Create(ctx context.Context, ns *composeTypes.Namespace) (*composeTypes.Namespace, error)
	}

	mappingReader interface {
		FindByID(ctx context.Context, id uint64) (*types.DmlMapping, error)
	}

	Applier struct {
		mapping   mappingReader
		moduleSvc ComposeModuleSvc
		nsSvc     ComposeNamespaceSvc
	}
)

func NewApplier(m mappingReader, mod ComposeModuleSvc, ns ComposeNamespaceSvc) *Applier {
	return &Applier{mapping: m, moduleSvc: mod, nsSvc: ns}
}

// Apply creates or updates the compose namespace + module for a single DmlMapping.
// Idempotent: existing module (matched by handle) is updated in-place.
func (a *Applier) Apply(ctx context.Context, mappingID uint64) error {
	mp, err := a.mapping.FindByID(ctx, mappingID)
	if err != nil {
		return err
	}
	if mp.Skip {
		return nil
	}

	ns, err := a.ensureNamespace(ctx, mp)
	if err != nil {
		return fmt.Errorf("dml apply: namespace: %w", err)
	}

	return a.applyMapping(ctx, ns, mp)
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

func (a *Applier) applyMapping(ctx context.Context, ns *composeTypes.Namespace, mp *types.DmlMapping) error {
	fields := make(composeTypes.ModuleFieldSet, 0, len(mp.Columns))
	for _, col := range mp.Columns {
		if col.Skip {
			continue
		}
		label := col.Label
		if label == "" {
			label = col.FieldName
		}
		fields = append(fields, &composeTypes.ModuleField{
			NamespaceID: ns.ID,
			Name:        col.FieldName,
			Label:       label,
			Kind:        col.FieldKind,
		})
	}

	existing, err := a.moduleSvc.FindByHandle(ctx, ns.ID, mp.ModuleHandle)
	if err == nil {
		existing.Fields = fields
		existing.Name = mp.ModuleName
		if existing.Name == "" {
			existing.Name = mp.ModuleHandle
		}
		_, err = a.moduleSvc.Update(ctx, existing)
		return err
	}

	name := mp.ModuleName
	if name == "" {
		name = mp.ModuleHandle
	}
	_, err = a.moduleSvc.Create(ctx, &composeTypes.Module{
		Handle:      mp.ModuleHandle,
		Name:        name,
		NamespaceID: ns.ID,
		Fields:      fields,
	})
	return err
}
