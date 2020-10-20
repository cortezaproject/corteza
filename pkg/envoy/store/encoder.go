package store

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/store"
	"time"
)

type (
	StoreEncoder struct {
		s store.Storer
	}
)

var (
	rvSanitizer = values.Sanitizer()
	rvValidator = values.Validator()
)

func NewStoreEncoder(s store.Storer) *StoreEncoder {
	return &StoreEncoder{
		s: s,
	}
}

func (se *StoreEncoder) Encode(ctx context.Context, s store.Storer, nn ...envoy.Node) error {
	return store.Tx(ctx, se.s, func(ctx context.Context, s store.Storer) error {
		for _, n := range nn {
			switch n.Resource() {
			case types.NamespaceRBACResource.String():
				ns := n.(*envoy.ComposeNamespaceNode)
				_, err := se.encodeNamespace(ctx, s, ns)
				if err != nil {
					return err
				}

			case types.ModuleRBACResource.String():
				mod := n.(*envoy.ComposeModuleNode)
				_, err := se.encodeModule(ctx, s, mod)
				if err != nil {
					return err
				}

			case "compose:record:":
				rec := n.(*envoy.ComposeRecordNode)
				err := se.encodeRecord(ctx, s, rec)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (se *StoreEncoder) encodeNamespace(ctx context.Context, s store.Storer, n *envoy.ComposeNamespaceNode) (*types.Namespace, error) {
	cns := n.Ns

	cns.ID = nextID()
	if cns.CreatedAt.IsZero() {
		cns.CreatedAt = time.Now()
	}

	if err := store.CreateComposeNamespace(ctx, s, cns); err != nil {
		return nil, err
	}

	return cns, nil
}

func (se *StoreEncoder) encodeModule(ctx context.Context, s store.Storer, m *envoy.ComposeModuleNode) (*types.Module, error) {
	mod := m.Mod
	cns := m.Ns

	mod.ID = nextID()
	mod.NamespaceID = cns.ID

	if mod.CreatedAt.IsZero() {
		mod.CreatedAt = time.Now()
	}

	// Store the module
	if err := store.CreateComposeModule(ctx, s, mod); err != nil {
		return nil, err
	}

	// Store module fields
	for _, f := range mod.Fields {
		f.ID = nextID()
		f.ModuleID = mod.ID

		if err := store.CreateComposeModuleField(ctx, s, f); err != nil {
			return nil, err
		}

		// Update the original module fields so dependant resources can proceed without issues
		mod.Fields = append(mod.Fields, f)
	}

	return mod, nil
}

func (se *StoreEncoder) encodeRecord(ctx context.Context, s store.Storer, m *envoy.ComposeRecordNode) error {
	mod := m.Mod

	return m.Walk(func(rec *types.Record) error {
		rec.ID = nextID()
		rec.ModuleID = mod.ID
		rec.NamespaceID = mod.NamespaceID

		if rec.CreatedAt.IsZero() {
			rec.CreatedAt = time.Now()
		}

		rec.Values = make(types.RecordValueSet, 0, 100)

		// Process record values
		for _, crv := range rec.Values {
			crv.RecordID = rec.ID
			rec.Values = append(rec.Values, crv)
		}
		rec.Values.SetUpdatedFlag(true)
		rec.Values = se.setDefaultComposeRecordValues(mod, rec.Values)
		rec.Values = rvSanitizer.Run(mod, rec.Values)
		err := rvValidator.Run(ctx, s, mod, rec)
		if err != nil {
			return err
		}

		if err := store.CreateComposeRecord(ctx, s, m.Mod, rec); err != nil {
			return err
		}

		return nil
	})
}

// @note this method is coppied over from the compose/service/record.
// Would it be better to unify the two methods?
func (se *StoreEncoder) setDefaultComposeRecordValues(m *types.Module, vv types.RecordValueSet) (out types.RecordValueSet) {
	out = vv

	for _, f := range m.Fields {
		if f.DefaultValue == nil {
			continue
		}

		for i, dv := range f.DefaultValue {
			// Default values on field are (might be) without field name and place
			if !out.Has(f.Name, uint(i)) {
				out = append(out, &types.RecordValue{
					Name:  f.Name,
					Value: dv.Value,
					Place: uint(i),
				})
			}
		}
	}

	return
}
