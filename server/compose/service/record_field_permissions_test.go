package service

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/stretchr/testify/require"
)

// denyFieldAC denies CanReadRecordValueOnModuleField for one field name and
// allows everything else.
type denyFieldAC struct{ denied string }

func (ac denyFieldAC) CanReadRecordValueOnModuleField(_ context.Context, f *types.ModuleField) bool {
	return f == nil || f.Name != ac.denied
}

func (ac denyFieldAC) CanUpdateRecordValueOnModuleField(context.Context, *types.ModuleField) bool {
	return true
}

func testModuleWithSecret() *types.Module {
	return &types.Module{Fields: types.ModuleFieldSet{
		{Name: "Title", Kind: "String"},
		{Name: "Secret", Kind: "String"},
	}}
}

// field identifiers resolve case-insensitively in the DAL, so a filter, sort or
// summary over a denied field must be refused however it is capitalised
func TestValidateFilterFieldPermissions_caseInsensitive(t *testing.T) {
	var (
		req = require.New(t)
		svc = record{}
		ac  = denyFieldAC{denied: "Secret"}
		m   = testModuleWithSecret()
	)

	for _, q := range []string{
		"Secret = 'x'",
		"secret = 'x'",
		"SECRET = 'x'",
		"Title = 'a' AND secret LIKE 'p%'",
	} {
		req.Error(svc.validateFilterFieldPermissions(context.Background(), ac, m, q), "query %q", q)
	}

	for _, q := range []string{"Title = 'a'", "title LIKE 'A%'", ""} {
		req.NoError(svc.validateFilterFieldPermissions(context.Background(), ac, m, q), "query %q", q)
	}
}

func TestValidateSortFieldPermissions_caseInsensitive(t *testing.T) {
	var (
		req = require.New(t)
		svc = record{}
		ac  = denyFieldAC{denied: "Secret"}
		m   = testModuleWithSecret()
	)

	for _, s := range []string{"Secret", "secret DESC", "Title, secret ASC"} {
		var sort filter.SortExprSet
		req.NoError(sort.Set(s))
		req.Error(svc.validateSortFieldPermissions(context.Background(), ac, m, sort), "sort %q", s)
	}

	for _, s := range []string{"Title", "title DESC"} {
		var sort filter.SortExprSet
		req.NoError(sort.Set(s))
		req.NoError(svc.validateSortFieldPermissions(context.Background(), ac, m, sort), "sort %q", s)
	}
}

func TestValidateSummaryFieldPermissions_caseInsensitive(t *testing.T) {
	var (
		req = require.New(t)
		svc = record{}
		ac  = denyFieldAC{denied: "Secret"}
		m   = testModuleWithSecret()
	)

	for _, f := range []string{"Secret", "secret", "SECRET"} {
		ss := []types.RecordSummaryReq{{Name: "min", Field: f}}
		req.Error(svc.validateSummaryFieldPermissions(context.Background(), ac, m, ss), "field %q", f)
	}

	ss := []types.RecordSummaryReq{{Name: "min", Field: "Title"}}
	req.NoError(svc.validateSummaryFieldPermissions(context.Background(), ac, m, ss))
}
