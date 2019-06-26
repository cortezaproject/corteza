package service

import (
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/internal/test"
)

func TestValueSanitizer(t *testing.T) {
	var (
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
	test.NoError(t, err, "unexpected error for sanitizeValues() call: %v")
	test.Assert(t, len(out) == 1, "expecting 1 record value after sanitization, got %d", len(rvs))

	rvs = types.RecordValueSet{{Name: "unknown", Value: "single"}}
	out, err = svc.sanitizeValues(module, rvs)
	test.Assert(t, err != nil, "expecting sanitizeValues() to return an error, got nil")

	rvs = types.RecordValueSet{{Name: "single1", Value: "single"}, {Name: "single1", Value: "single2"}}
	out, err = svc.sanitizeValues(module, rvs)
	test.Assert(t, err != nil, "expecting sanitizeValues() to return an error, got nil")

	rvs = types.RecordValueSet{{Name: "multi1", Value: "multi1"}, {Name: "multi1", Value: "multi1"}}
	out, err = svc.sanitizeValues(module, rvs)
	test.NoError(t, err, "unexpected error for sanitizeValues() call: %v")
	test.Assert(t, len(out) == 2, "expecting 2 record values after sanitization, got %d", len(rvs))
	test.Assert(t, out[0].Place == 0, "expecting first value to have place value 0, got %d", out[0].Place)
	test.Assert(t, out[1].Place == 1, "expecting second value to have place value 1, got %d", out[1].Place)

	rvs = types.RecordValueSet{{Name: "ref1", Value: "multi1"}}
	out, err = svc.sanitizeValues(module, rvs)
	test.Assert(t, err != nil, "expecting sanitizeValues() to return an error, got nil")

	rvs = types.RecordValueSet{{Name: "ref1", Value: "12345"}}
	out, err = svc.sanitizeValues(module, rvs)
	test.NoError(t, err, "unexpected error for sanitizeValues() call: %v")
	test.Assert(t, len(out) == 1, "expecting 1 record values after sanitization, got %d", len(rvs))
	test.Assert(t, out[0].Ref == 12345, "expecting parsed ref value to match, got %d", rvs[0].Ref)

	rvs = types.RecordValueSet{{Name: "multiRef1", Value: "12345"}, {Name: "multiRef1", Value: "67890"}}
	out, err = svc.sanitizeValues(module, rvs)
	test.NoError(t, err, "unexpected error for sanitizeValues() call: %v")
	test.Assert(t, len(rvs) == 2, "expecting 2 record values after sanitization, got %d", len(rvs))
	test.Assert(t, out[0].Ref == 12345, "expecting parsed ref value to match, got %d", out[0].Ref)
	test.Assert(t, out[1].Ref == 67890, "expecting parsed ref value to match, got %d", out[1].Ref)

	rvs = types.RecordValueSet{{Name: "ref1", Value: ""}}
	out, err = svc.sanitizeValues(module, rvs)
	test.NoError(t, err, "unexpected error for sanitizeValues() call: %v")
	test.Assert(t, len(out) == 0, "expecting 0 record values after sanitization, got %d", len(rvs))
}
