package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

func TestGeneralValueSetValidation(t *testing.T) {
	var (
		req = require.New(t)

		svc = record{
			ctx: context.Background(),
			ac:  AccessControl(&permissions.ServiceAllowAll{}),
		}
		module = &types.Module{
			Fields: types.ModuleFieldSet{
				&types.ModuleField{Name: "single1"},
				&types.ModuleField{Name: "multi1", Multi: true},
				&types.ModuleField{Name: "ref1", Kind: "Record"},
				&types.ModuleField{Name: "multiRef1", Kind: "Record", Multi: true},
			},
		}

		rvs types.RecordValueSet
		err error
	)

	rvs = types.RecordValueSet{{Name: "single1", Value: "single"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.NoError(err)

	rvs = types.RecordValueSet{{Name: "unknown", Value: "single"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.True(err != nil, "expecting generalValueSetValidation() to return an error, got nil")

	rvs = types.RecordValueSet{{Name: "single1", Value: "single"}, {Name: "single1", Value: "single2"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.Error(err)

	rvs = types.RecordValueSet{{Name: "multi1", Value: "multi1"}, {Name: "multi1", Value: "multi1"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.NoError(err)

	rvs = types.RecordValueSet{{Name: "ref1", Value: "non numeric value"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.Error(err)

	rvs = types.RecordValueSet{{Name: "ref1", Value: "12345"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.NoError(err)

	rvs = types.RecordValueSet{{Name: "multiRef1", Value: "12345"}, {Name: "multiRef1", Value: "67890"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.NoError(err)
	req.Len(rvs, 2, "expecting 2 record values after sanitization, got %d", len(rvs))

	rvs = types.RecordValueSet{{Name: "ref1", Value: ""}}
	err = svc.generalValueSetValidation(module, rvs)
	req.NoError(err)
}

func TestDefaultValueSetting(t *testing.T) {
	var (
		a = assert.New(t)

		svc = record{
			ac: AccessControl(&permissions.ServiceAllowAll{}),
		}
		mod = &types.Module{
			Fields: types.ModuleFieldSet{
				&types.ModuleField{Name: "single", DefaultValue: []*types.RecordValue{{Value: "s"}}},
				&types.ModuleField{Name: "multi", Multi: true, DefaultValue: []*types.RecordValue{{Value: "m1", Place: 0}, {Value: "m2", Place: 1}}},
			},
		}

		chk = func(vv types.RecordValueSet, f string, p uint, v string) {
			a.True(vv.Has("multi", p))
			a.Equal(v, vv.Get(f, p).Value)
		}
	)

	out := svc.setDefaultValues(mod, nil)
	chk(out, "single", 0, "s")
	chk(out, "multi", 0, "m1")
	chk(out, "multi", 1, "m2")
}

func TestProcUpdateOwnerPreservation(t *testing.T) {
	var (
		a = assert.New(t)

		svc = record{
			sanitizer: values.Sanitizer(),
			validator: values.Validator(),
		}

		mod = &types.Module{
			Fields: types.ModuleFieldSet{},
		}

		oldRec = &types.Record{
			OwnedBy: 1,
			Values:  types.RecordValueSet{},
		}
		newRec = &types.Record{
			OwnedBy: 0,
			Values:  types.RecordValueSet{},
		}
	)

	svc.procUpdate(10, mod, newRec, oldRec)
	a.Equal(newRec.OwnedBy, uint64(1))
	svc.procUpdate(10, mod, newRec, oldRec)
	a.Equal(newRec.OwnedBy, uint64(1))
}

func TestProcUpdateOwnerChanged(t *testing.T) {
	var (
		a = assert.New(t)

		svc = record{
			sanitizer: values.Sanitizer(),
			validator: values.Validator(),
		}

		mod = &types.Module{
			Fields: types.ModuleFieldSet{},
		}

		oldRec = &types.Record{
			OwnedBy: 1,
			Values:  types.RecordValueSet{},
		}
		newRec = &types.Record{
			OwnedBy: 9,
			Values:  types.RecordValueSet{},
		}
	)

	svc.procUpdate(10, mod, newRec, oldRec)
	a.Equal(newRec.OwnedBy, uint64(9))
	svc.procUpdate(10, mod, newRec, oldRec)
	a.Equal(newRec.OwnedBy, uint64(9))
}
