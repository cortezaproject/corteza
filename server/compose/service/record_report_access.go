package service

import (
	"context"
	"strconv"
	"strings"
	"sync"

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
	denying, ok := recordReadDenyingRoles(ac, m)
	if !ok {
		return true
	}
	if len(denying) == 0 {
		return false
	}

	// nil when the roles cannot be listed: every deny counts
	var applicable map[uint64]bool
	if src, ok := ac.(interface {
		applicableRoles(context.Context) (map[uint64]bool, bool)
	}); ok {
		applicable, _ = src.applicableRoles(ctx)
	}
	if applicable == nil {
		return true
	}

	for role := range denying {
		if applicable[role] {
			return true
		}
	}
	return false
}

// recordReadDenyingRoles lists the roles a rule denies reading the module's
// records to, and whether the rules could be listed
func recordReadDenyingRoles(ac any, m *types.Module) (map[uint64]bool, bool) {
	if src, ok := ac.(interface {
		recordReadDenies() (recordReadDenyIndex, bool)
	}); ok {
		idx, ok := src.recordReadDenies()
		return idx.lookup(m), ok
	}

	src, ok := ac.(interface{ ruleSet() (rbac.RuleSet, bool) })
	if !ok {
		return nil, false
	}
	rr, ok := src.ruleSet()
	if !ok {
		return nil, false
	}

	return indexRecordReadDenies(rr).lookup(m), true
}

// recordReadDenyIndex holds the roles denied reading records, by
// "namespace/module" (either part may be "*")
type recordReadDenyIndex map[string]map[uint64]bool

func indexRecordReadDenies(rr rbac.RuleSet) recordReadDenyIndex {
	var (
		idx    = recordReadDenyIndex{}
		prefix = types.RecordResourceType + "/"
	)

	for _, r := range rr {
		if r.Operation != "read" || r.Access != rbac.Deny || !strings.HasPrefix(r.Resource, prefix) {
			continue
		}

		path := strings.Split(strings.TrimPrefix(r.Resource, prefix), "/")
		if len(path) != 3 {
			continue
		}

		key := path[0] + "/" + path[1]
		if idx[key] == nil {
			idx[key] = make(map[uint64]bool)
		}
		idx[key][r.RoleID] = true
	}

	return idx
}

// lookup lists the roles denied reading the module's records: by the module,
// its namespace or every namespace and module
func (idx recordReadDenyIndex) lookup(m *types.Module) map[uint64]bool {
	var (
		ns  = strconv.FormatUint(m.NamespaceID, 10)
		mod = strconv.FormatUint(m.ID, 10)
		out = make(map[uint64]bool)
	)

	for _, key := range []string{ns + "/" + mod, ns + "/*", "*/" + mod, "*/*"} {
		for role := range idx[key] {
			out[role] = true
		}
	}

	return out
}

// readDenies keeps the record read deny index of the RBAC rules, rebuilt
// only when the rules change, so list and report requests do not scan every
// rule
var readDenies struct {
	mux     sync.Mutex
	src     any
	version uint64
	idx     recordReadDenyIndex
}

// recordReadDenies returns the record read deny index, and whether the RBAC
// service could provide the rules for it
func (svc accessControl) recordReadDenies() (recordReadDenyIndex, bool) {
	src, ok := svc.rbac.(interface {
		RulesVersion() uint64
		Rules() rbac.RuleSet
	})
	if !ok {
		return nil, false
	}

	// version first: rules read after it are at least as new
	version := src.RulesVersion()

	readDenies.mux.Lock()
	defer readDenies.mux.Unlock()

	if readDenies.src != src || readDenies.version != version || readDenies.idx == nil {
		readDenies.idx = indexRecordReadDenies(src.Rules())
		readDenies.src, readDenies.version = src, version
	}

	return readDenies.idx, true
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
