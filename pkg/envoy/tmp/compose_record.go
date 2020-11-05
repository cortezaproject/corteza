package tmp

import (
	"context"
	"errors"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composeRecordPreproc struct {
		es *encoderState
		s  store.Storer
	}
)

func NewComposeRecordPreproc(is *encoderState, s store.Storer) *composeRecordPreproc {
	return &composeRecordPreproc{
		es: is,
		s:  s,
	}
}

func (p *composeRecordPreproc) Process(ctx context.Context, state *envoy.ExecState) error {
	res, is := state.Res.(*resource.ComposeRecord)
	if !is {
		return nil
	}

	// Get relate namespace
	ns, err := findNamespace(ctx, p.s, state.ParentResources, res.NsRef.Identifiers)
	if err != nil {
		return err
	}
	if ns == nil {
		return errors.New("@todo couldn't resolve namespace")
	}
	p.es.Set(res, res.NsRef.ResourceType, ns.ID, res.NsRef.Identifiers.StringSlice()...)

	// Get relate namespace
	mod, err := findModule(ctx, p.s, state.ParentResources, res.ModRef.Identifiers)
	if err != nil {
		return err
	}
	if mod == nil {
		return errors.New("@todo couldn't resolve module")
	}
	p.es.Set(res, res.ModRef.ResourceType, mod.ID, res.ModRef.Identifiers.StringSlice()...)

	// @todo existing records
	//
	// Hookup with labels to determine existing records

	ff, err := loadComposeModuleFields(ctx, p.s, mod)
	if err != nil {
		return err
	}
	res.ModFields = ff

	for _, f := range ff {
		switch f.Kind {
		case "Record":
			refM := f.Options.String("module")
			if refM != "" && refM != "0" {
				// Make a reference with that module's records
				res.AddRef(resource.COMPOSE_RECORD_RESOURCE_TYPE, refM)
			}
		}
	}

	return nil
}

func encodeComposeRecord(ctx context.Context, ectx *encodingContext, s store.Storer, state resRefs, res *resource.ComposeRecord) (resRefs, error) {
	var err error
	rState := make(resRefs)

	// Namespace...
	nsID := uint64(0)
	for _, v := range state[resource.COMPOSE_NAMESPACE_RESOURCE_TYPE] {
		nsID = v
		break
	}

	// Module...
	modID := uint64(0)
	for _, v := range state[resource.COMPOSE_MODULE_RESOURCE_TYPE] {
		modID = v
		break
	}
	mod := &types.Module{
		NamespaceID: nsID,
		ID:          modID,
		Fields:      res.ModFields,
	}

	return rState, res.Walker(func(r *resource.ComposeRecordRaw) error {
		rec := &types.Record{
			NamespaceID: nsID,
			ModuleID:    modID,
		}

		rec.ID = state[res.ResourceType()][r.ID]
		if rec.ID <= 0 {
			rec.ID = nextID()
		}
		rState.Set(resource.COMPOSE_RECORD_RESOURCE_TYPE, rec.ID, r.ID)

		for k, v := range r.SysValues {
			if v == "" {
				continue
			}

			switch k {
			case "createdAt":
				// @todo set time
				rec.CreatedAt = time.Now()

			case "updatedAt":
				// @todo set time
				rec.UpdatedAt = nil

			case "deletedAt":
				// @todo set time
				rec.DeletedAt = nil
			}
		}

		rvs := make(types.RecordValueSet, 0, len(r.Values))
		for k, v := range r.Values {
			rv := &types.RecordValue{
				RecordID: rec.ID,
				Name:     k,
				Value:    v,
				Updated:  true,
			}

			rvs = append(rvs, rv)
		}

		rec.Values = rvSanitizer.Run(mod, rvs)

		if !ectx.partial && !ectx.exists {
			err = store.CreateComposeRecord(ctx, s, mod, rec)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
