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
	if readsAllRecords(ctx, svc.ac, m) {
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

// readsAllRecords answers whether the caller reads every record of the module,
// from the rules alone: a record is hidden only by a rule denying it, so
// whoever reads a record no rule names, through roles that apply to every
// record alike, with no such deny for a role they can hold, reads them all.
func readsAllRecords(ctx context.Context, ac any, m *types.Module) bool {
	// The probe needs an ID: a wildcard resource never passes a check. No rule
	// names this one, and nobody owns the record.
	probe := &types.Record{ID: ^uint64(0), NamespaceID: m.NamespaceID, ModuleID: m.ID}
	probe.SetModule(m)

	// Contextual roles are left out: one whose expression happens to match the
	// ownerless probe says nothing about the records that are owned
	src, ok := ac.(interface {
		readsWithoutContextRoles(context.Context, *types.Record) (bool, bool)
	})
	if !ok {
		return false
	}
	if can, ok := src.readsWithoutContextRoles(ctx, probe); !ok || !can {
		return false
	}

	return !recordReadDenied(ctx, ac, m)
}

// recordReadDenied answers whether any rule that can apply to the caller
// denies reading the module's records, or when the rules cannot be listed,
// whether one may.
//
// A deny given to a role the caller can not hold hides nothing from them;
// counting it would send everybody's reports through the records one by one.
func recordReadDenied(ctx context.Context, ac any, m *types.Module) bool {
	src, ok := ac.(interface{ ruleSet() (rbac.RuleSet, bool) })
	if !ok {
		return true
	}
	rr, ok := src.ruleSet()
	if !ok {
		return true
	}

	// nil when the roles cannot be listed: every deny counts
	var applicable map[uint64]bool
	if src, ok := ac.(interface {
		applicableRoles(context.Context) (map[uint64]bool, bool)
	}); ok {
		applicable, _ = src.applicableRoles(ctx)
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
		if applicable != nil && !applicable[r.RoleID] {
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

// readsWithoutContextRoles answers whether the caller reads the record through
// roles other than contextual ones, and whether the RBAC service could tell.
func (svc accessControl) readsWithoutContextRoles(ctx context.Context, r *types.Record) (bool, bool) {
	src, ok := svc.rbac.(interface {
		CheckWithoutContextRoles(rbac.Session, string, rbac.Resource) rbac.Access
	})
	if !ok {
		return false, false
	}
	return src.CheckWithoutContextRoles(rbac.ContextToSession(ctx), "read", r) == rbac.Allow, true
}

// applicableRoles lists the roles whose rules can apply to the caller, and
// whether the RBAC service could.
func (svc accessControl) applicableRoles(ctx context.Context) (map[uint64]bool, bool) {
	src, ok := svc.rbac.(interface {
		ApplicableRoles(rbac.Session) map[uint64]bool
	})
	if !ok {
		return nil, false
	}
	return src.ApplicableRoles(rbac.ContextToSession(ctx)), true
}
