package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza/server/compose/dalutils"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
)

// reportReadableCap is how many records a report reads one by one for a caller
// who may read only some of the module's records.
const reportReadableCap = 10000

// reportReadable answers which of the module's records a report may aggregate:
// all of them (nil), or the IDs of the ones the caller may read.
//
// A report adds records up in the database, past the per-record check a list
// runs. Whoever reads a record nobody owns, on a module whose records no rule
// denies reading, reads them all and keeps that; anybody else, an owner reading
// their own records for one, has the readable ones collected first.
func (svc record) reportReadable(ctx context.Context, m *types.Module, f string) (ids []any, err error) {
	// The probe needs an ID: a wildcard resource never passes a check. No rule
	// names this one, and nobody owns the record.
	probe := &types.Record{ID: ^uint64(0), NamespaceID: m.NamespaceID, ModuleID: m.ID}
	probe.SetModule(m)
	if svc.ac.CanReadRecord(ctx, probe) && !recordReadDenied(svc.ac, m) {
		return nil, nil
	}

	set, _, err := dalutils.ComposeRecordsList(ctx, svc.dal, m, types.RecordFilter{
		NamespaceID: m.NamespaceID,
		ModuleID:    m.ID,
		Query:       f,
		Deleted:     filter.StateExcluded,
		Paging:      filter.Paging{Limit: reportReadableCap + 1},
		Check:       ComposeRecordFilterChecker(ctx, svc.ac, m),
	})
	if err != nil {
		return nil, err
	}
	if len(set) > reportReadableCap {
		return nil, RecordErrReportTooManyReadable()
	}

	ids = make([]any, 0, len(set))
	for _, r := range set {
		ids = append(ids, r.ID)
	}
	return ids, nil
}

// recordReadDenied answers whether any rule denies reading the module's
// records, or when the rules cannot be listed, whether one may.
func recordReadDenied(ac any, m *types.Module) bool {
	src, ok := ac.(interface{ ruleSet() (rbac.RuleSet, bool) })
	if !ok {
		return true
	}
	rr, ok := src.ruleSet()
	if !ok {
		return true
	}

	var (
		prefix = types.RecordResourceType + "/"
		ns     = strconv.FormatUint(m.NamespaceID, 10)
		mod    = strconv.FormatUint(m.ID, 10)
	)
	for _, r := range rr {
		if r.Operation != "read" || r.Access != rbac.Deny || !strings.HasPrefix(r.Resource, prefix) {
			continue
		}
		path := strings.Split(strings.TrimPrefix(r.Resource, prefix), "/")
		if len(path) == 3 && (path[0] == "*" || path[0] == ns) && (path[1] == "*" || path[1] == mod) {
			return true
		}
	}
	return false
}

// ruleSet lists every rule, and whether the RBAC service could.
func (svc accessControl) ruleSet() (rbac.RuleSet, bool) {
	if src, ok := svc.rbac.(interface{ Rules() rbac.RuleSet }); ok {
		return src.Rules(), true
	}
	return nil, false
}
