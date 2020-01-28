package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

func TestValueSanitizer(t *testing.T) {
	var (
		req = require.New(t)

		svc = record{
			ac: AccessControl(&permissions.ServiceAllowAll{}),
		}
		module = &types.Module{
			Fields: types.ModuleFieldSet{
				&types.ModuleField{Name: "single1"},
				&types.ModuleField{Name: "multi1", Multi: true},
				&types.ModuleField{Name: "ref1", Kind: "Record"},
				&types.ModuleField{Name: "multiRef1", Kind: "Record", Multi: true},
			},
		}

		rvs, out types.RecordValueSet
		err      error
	)

	rvs = types.RecordValueSet{{Name: "single1", Value: "single"}}
	out, err = svc.sanitizeValues(module, rvs)
	req.NoError(err)
	req.Len(out, 1, "expecting 1 record value after sanitization, got %d", len(rvs))

	rvs = types.RecordValueSet{{Name: "unknown", Value: "single"}}
	out, err = svc.sanitizeValues(module, rvs)
	req.True(err != nil, "expecting sanitizeValues() to return an error, got nil")

	rvs = types.RecordValueSet{{Name: "single1", Value: "single"}, {Name: "single1", Value: "single2"}}
	out, err = svc.sanitizeValues(module, rvs)
	req.True(err != nil, "expecting sanitizeValues() to return an error, got nil")

	rvs = types.RecordValueSet{{Name: "multi1", Value: "multi1"}, {Name: "multi1", Value: "multi1"}}
	out, err = svc.sanitizeValues(module, rvs)
	req.NoError(err)
	req.Len(out, 2, "expecting 2 record values after sanitization, got %d", len(rvs))
	req.Equal(uint(0), out[0].Place, "expecting first value to have place value 0, got %d", out[0].Place)
	req.Equal(uint(1), out[1].Place, "expecting second value to have place value 1, got %d", out[1].Place)

	rvs = types.RecordValueSet{{Name: "ref1", Value: "multi1"}}
	out, err = svc.sanitizeValues(module, rvs)
	req.True(err != nil, "expecting sanitizeValues() to return an error, got nil")

	rvs = types.RecordValueSet{{Name: "ref1", Value: "12345"}}
	out, err = svc.sanitizeValues(module, rvs)
	req.NoError(err)
	req.Len(out, 1, "expecting 1 record values after sanitization, got %d", len(rvs))
	req.Equal(uint64(12345), out[0].Ref, "expecting parsed ref value to match, got %d", rvs[0].Ref)

	rvs = types.RecordValueSet{{Name: "multiRef1", Value: "12345"}, {Name: "multiRef1", Value: "67890"}}
	out, err = svc.sanitizeValues(module, rvs)
	req.NoError(err)
	req.Len(rvs, 2, "expecting 2 record values after sanitization, got %d", len(rvs))
	req.Equal(uint64(12345), out[0].Ref, "expecting parsed ref value to match, got %d", out[0].Ref)
	req.Equal(uint64(67890), out[1].Ref, "expecting parsed ref value to match, got %d", out[1].Ref)

	rvs = types.RecordValueSet{{Name: "ref1", Value: ""}}
	out, err = svc.sanitizeValues(module, rvs)
	req.NoError(err)
	req.Len(out, 0, "expecting 0 record values after sanitization, got %d", len(rvs))
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

		rec *types.Record

		chk = func(vv types.RecordValueSet, f string, p uint, v string) {
			a.True(vv.Has("multi", p))
			a.Equal(v, vv.Get(f, p).Value)
		}
	)

	rec = &types.Record{}
	a.NoError(svc.setDefaultValues(mod, rec))
	chk(rec.Values, "single", 0, "s")
	chk(rec.Values, "multi", 0, "m1")
	chk(rec.Values, "multi", 1, "m2")

}
