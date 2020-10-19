package encoder

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/compose/service/values"
	compTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/types"
	"github.com/cortezaproject/corteza-server/store"
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

func (se *StoreEncoder) Encode(ctx context.Context, nn ...types.Node) error {
	for _, n := range nn {
		switch n.Resource() {
		case compTypes.NamespaceRBACResource.String():
			ns := n.(*types.ComposeNamespaceNode)
			_, err := se.encodeNamespace(ctx, ns)
			if err != nil {
				return err
			}

		case compTypes.ModuleRBACResource.String():
			mod := n.(*types.ComposeModuleNode)
			_, err := se.encodeModule(ctx, mod)
			if err != nil {
				return err
			}

		case "compose:record:":
			rec := n.(*types.ComposeRecordNode)
			err := se.encodeRecord(ctx, rec)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (se *StoreEncoder) encodeNamespace(ctx context.Context, n *types.ComposeNamespaceNode) (*types.ComposeNamespace, error) {
	cns := n.Ns

	err := store.Tx(ctx, se.s, func(ctx context.Context, s store.Storer) error {
		cns.ID = nextID()
		if cns.CreatedAt.IsZero() {
			cns.CreatedAt = time.Now()
		}

		if err := store.CreateComposeNamespace(ctx, s, &cns.Namespace); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return cns, nil
}

func (se *StoreEncoder) encodeModule(ctx context.Context, m *types.ComposeModuleNode) (*types.ComposeModule, error) {
	cmod := m.Mod
	cns := m.Ns

	err := store.Tx(ctx, se.s, func(ctx context.Context, s store.Storer) error {
		cmod.ID = nextID()
		cmod.NamespaceID = cns.ID

		if cmod.CreatedAt.IsZero() {
			cmod.CreatedAt = time.Now()
		}

		// Store the module
		mod := &cmod.Module
		if err := store.CreateComposeModule(ctx, s, mod); err != nil {
			return err
		}

		// Store module fields
		for _, cf := range cmod.Fields {
			cf.ID = nextID()
			cf.ModuleID = cmod.ID

			f := &cf.ModuleField
			if err := store.CreateComposeModuleField(ctx, s, f); err != nil {
				return err
			}

			// Update the original module fields so dependant resourcess can proceed without issues
			mod.Fields = append(mod.Fields, f)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return cmod, nil
}

func (se *StoreEncoder) encodeRecord(ctx context.Context, m *types.ComposeRecordNode) error {
	cmod := m.Mod
	err := store.Tx(ctx, se.s, func(ctx context.Context, s store.Storer) error {
		return m.Walk(func(cr *types.ComposeRecord) error {
			cr.ID = nextID()
			cr.ModuleID = cmod.ID
			cr.NamespaceID = cmod.NamespaceID

			if cr.CreatedAt.IsZero() {
				cr.CreatedAt = time.Now()
			}

			rec := &cr.Record
			rec.Values = make(compTypes.RecordValueSet, 0, 100)

			// Process record values
			for _, crv := range cr.Values {
				crv.RecordID = rec.ID
				rec.Values = append(rec.Values, &crv.RecordValue)
			}
			rec.Values.SetUpdatedFlag(true)
			rec.Values = se.setDefaultComposeRecordValues(cmod, rec.Values)
			rec.Values = rvSanitizer.Run(&cmod.Module, rec.Values)
			err := rvValidator.Run(ctx, s, &cmod.Module, rec)
			if err != nil {
				return err
			}

			if err := store.CreateComposeRecord(ctx, s, &m.Mod.Module, rec); err != nil {
				return err
			}

			return nil
		})
	})

	if err != nil {
		return err
	}

	return nil
}

// @note this method is coppied over from the compose/service/record.
// Would it be better to unify the two methods?
func (se *StoreEncoder) setDefaultComposeRecordValues(m *types.ComposeModule, vv compTypes.RecordValueSet) (out compTypes.RecordValueSet) {
	out = vv

	for _, f := range m.Fields {
		if f.DefaultValue == nil {
			continue
		}

		for i, dv := range f.DefaultValue {
			// Default values on field are (might be) without field name and place
			if !out.Has(f.Name, uint(i)) {
				out = append(out, &compTypes.RecordValue{
					Name:  f.Name,
					Value: dv.Value,
					Place: uint(i),
				})
			}
		}
	}

	return
}
