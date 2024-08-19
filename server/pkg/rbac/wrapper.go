package rbac

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/davecgh/go-spew/spew"
)

type (
	WrapperConfig struct {
		InitialIndexedRoles []uint64
		MaxIndexSize        int
	}

	wrapperService struct {
		cfg WrapperConfig

		store   rbacRulesStore
		counter *usageCounter
		index   *wrapperIndex
		roles   []*Role
	}
)

func dftWrapperCfg(base WrapperConfig) (out WrapperConfig) {
	out = base

	if base.MaxIndexSize == 0 {
		out.MaxIndexSize = -1
	}

	return out
}

func Wrapper(ctx context.Context, store rbacRulesStore, cc WrapperConfig) (x *wrapperService, err error) {
	cc = dftWrapperCfg(cc)

	uc := &usageCounter{
		incChan: make(chan uint64, 256),
		sigChan: make(chan counterEntry, 8),
	}

	x = &wrapperService{
		cfg: cc,

		store:   store,
		counter: uc,
	}

	x.roles, err = x.loadRoles(ctx, store)
	if err != nil {
		return
	}

	x.index, err = x.loadIndex(ctx, store, x.roles)
	if err != nil {
		return
	}

	uc.watch(ctx)
	x.watch(ctx)

	return
}

func (svc *wrapperService) Clear() {
	svc.store = nil
	svc.counter = nil
	svc.index = nil
	svc.roles = nil
}

func (svc *wrapperService) Can(ses Session, op string, res Resource) (ok bool, err error) {
	ac, err := svc.Check(ses, op, res)
	if err != nil {
		return
	}

	return ac == Allow, nil
}

func (svc *wrapperService) Check(ses Session, op string, res Resource) (a Access, err error) {
	if hasWildcards(res.RbacResource()) {
		// prevent use of wildcard resources for checking permissions
		return Inherit, nil
	}

	fRoles := getContextRoles(ses, res, svc.roles...)

	return svc.check(ses.Context(), fRoles, op, res.RbacResource())
}

func (svc *wrapperService) check(ctx context.Context, rolesByKind partRoles, op, res string) (a Access, err error) {
	if member(rolesByKind, AnonymousRole) && len(rolesByKind) > 1 {
		// Integrity check; when user is member of anonymous role
		// should not be member of any other type of role
		return resolve(nil, Deny, failedIntegrityCheck), nil
	}

	if member(rolesByKind, BypassRole) {
		// if user has at least one bypass role, we allow access
		return resolve(nil, Allow, bypassRoleMembership), nil
	}

	// if indexedRules.empty() {
	// 	// no rules to check
	// 	return resolve(nil, Inherit, noRules)
	// }

	var (
		match   *Rule
		allowed bool
	)

	indexed, unindexed, err := svc.segmentRoles(ctx, rolesByKind)
	if err != nil {
		return Inherit, err
	}

	//
	// if trace != nil {
	// 	// from this point on, there is a chance trace (if set)
	// 	// will contain some rules.
	// 	//
	// 	// Stable order needs to be ensured: there is no production
	// 	// code that relies on that but tests might fail and API
	// 	// response would be flaky.
	// 	defer sortTraceRules(trace)
	// }

	st := evlState{
		op:  op,
		res: res,

		unindexedRoles: unindexed,
		indexedRoles:   indexed,
	}

	st.unindexedRules, err = svc.pullUnindexed(ctx, unindexed, op, res)
	if err != nil {
		return Inherit, err
	}

	// Priority is important here. We want to have
	// stable RBAC check behaviour and ability
	// to override allow/deny depending on how niche the role (type) is:
	//  - context (eg owners) are more niche than common
	//  - rules for common roles are more important than authenticated and anonymous role types
	//
	// Note that bypass roles are intentionally ignored here; if user is member of
	// bypass role there is no need to check any other rule
	for _, kind := range []roleKind{ContextRole, CommonRole, AuthenticatedRole, AnonymousRole} {
		// not a member of any role of this kind
		if len(rolesByKind[kind]) == 0 {
			continue
		}

		// reset allowed to false
		// for each role kind
		allowed = false

		for r := range rolesByKind[kind] {
			match = svc.getMatching(st, kind, r)

			// check all rules for each role the security-context
			if match == nil {
				// no rules match
				continue
			}

			// if trace != nil {
			// 	// if trace is enabled, append
			// 	// each matching rule
			// 	trace.Rules = append(trace.Rules, match)
			// }

			if match.Access == Deny {
				// if we stumble upon Deny we short-circuit the check
				return resolve(nil, Deny, ""), nil
			}

			if match.Access == Allow {
				// allow rule found, we need to check rules on other roles
				// before we allow it
				allowed = true
			}
		}

		if allowed {
			// at least one of the roles (per role type) in the security context
			// allows operation on a resource
			return resolve(nil, Allow, ""), nil
		}
	}

	// No rule matched
	return resolve(nil, Inherit, noMatch), nil
}

func (svc *wrapperService) segmentRoles(ctx context.Context, roles partRoles) (indexed, unindexed partRoles, err error) {
	unindexed = partRoles{}
	indexed = partRoles{}

	unindexed[CommonRole] = make(map[uint64]bool)
	indexed[CommonRole] = make(map[uint64]bool)

	for k, rg := range roles {
		for r := range rg {
			if svc.index.hasRole(r) {
				indexed[k][r] = true
				continue
			}

			unindexed[k][r] = true
		}
	}

	return
}

type (
	evlState struct {
		unindexedRoles partRoles
		indexedRoles   partRoles

		unindexedRules [5]map[uint64][]*Rule

		res string
		op  string
	}
)

func (svc *wrapperService) getMatching(st evlState, kind roleKind, role uint64) (rule *Rule) {
	var (
		aux   []*Rule
		rules RuleSet
	)

	// Indexed
	aux = svc.index.get(role, st.op, st.res)
	rules = append(rules, aux...)

	// Unindexed
	aux = st.unindexedRules[kind][role]
	rules = append(rules, aux...)

	set := RuleSet(rules)
	sort.Sort(set)

	for _, s := range set {
		if s.Access == Inherit {
			continue
		}

		return s
	}

	return nil
}

func (svc *wrapperService) pullUnindexed(ctx context.Context, unindexed partRoles, op, res string) (out [5]map[uint64][]*Rule, err error) {
	resPerm := make([]string, 0, 8)
	resPerm = append(resPerm, res)

	// Get all the resource permissions
	// @todo get permissions for parent resources; this will probs be some lookup table
	rr := strings.Split(res, "/")
	for i := len(rr) - 1; i >= 0; i-- {
		rr[i] = "*"
		resPerm = append(resPerm, strings.Join(rr, "/"))
	}

	for rk, rr := range unindexed {
		for r := range rr {
			auxRr := make([]*Rule, 0, 4)
			auxRr, _, err = svc.store.SearchRbacRules(ctx, RuleFilter{
				RoleID:    r,
				Resource:  resPerm,
				Operation: op,
			})
			if err != nil {
				return
			}

			if out[rk] == nil {
				out[rk] = map[uint64][]*Rule{
					r: auxRr,
				}
			} else {
				out[rk][r] = auxRr
			}
		}
	}

	return
}

func (svc *wrapperService) IndexRoleChange(ctx context.Context, roleID uint64) (err error) {
	aux, _, err := svc.store.SearchRbacRules(ctx, RuleFilter{
		RoleID: roleID,
	})
	if err != nil {
		return
	}

	// @todo cap this
	if len(svc.index.rules.children) > svc.cfg.MaxIndexSize {
		// @note probably remove a few extra just to avoid constantly doing this
		// @todo is this a good idea? Not sure if worth it since all of this is behind the scene anyways
		wp := svc.counter.worstPerformers(4)
		svc.index.remove(wp...)
	}

	svc.index.add(aux...)
	return
}

func (svc *wrapperService) watch(ctx context.Context) {
	t := time.NewTicker(time.Minute * 5)

	go func() {
		for {
			select {
			case <-t.C:
				spew.Dump("ticking")

			case change := <-svc.counter.sigChan:
				err := svc.IndexRoleChange(ctx, change.key)
				if err != nil {
					spew.Dump("wrapper watch change err", err)
				}

			case <-ctx.Done():
				return
			}
		}
	}()
}

// // // // // // // // // // // // // // // // // // // // // // // // // //

func makeKey(op, res string, role uint64) string {
	return fmt.Sprintf("%d:%s:%s", role, op, res)
}

//

// // // // // // // // // // // // // // // // // // // // // // // // // //
// Boilerplate & state management stuff

func (svc *wrapperService) loadRoles(ctx context.Context, s rbacRulesStore) (out []*Role, err error) {
	auxRoles, _, err := s.SearchRoles(ctx, types.RoleFilter{
		Paging: filter.Paging{
			Limit: 0,
		},
	})
	if err != nil {
		return
	}

	for _, ar := range auxRoles {
		out = append(out, &Role{
			id:     ar.ID,
			handle: ar.Handle,
			kind:   CommonRole,
		})
	}

	return
}

func (svc *wrapperService) loadIndex(ctx context.Context, s rbacRulesStore, allRoles []*Role) (out *wrapperIndex, err error) {
	// @todo smarter way to figure out what/how many roles we want to load up
	roles := svc.getIndexRoles(allRoles)

	rules := make(RuleSet, 0, 1024)
	var aux RuleSet
	for _, role := range roles {
		aux, _, err = s.SearchRbacRules(ctx, RuleFilter{
			RoleID: role.id,
			Limit:  0,
		})
		if err != nil {
			return
		}

		rules = append(rules, aux...)
	}

	out = &wrapperIndex{
		rules: buildRuleIndex(rules),
	}

	return
}

func (svc *wrapperService) getIndexRoles(allRoles []*Role) (out []*Role) {
	// User-specified what we want to index; respect that to the t
	if len(svc.cfg.InitialIndexedRoles) > 0 {
		for _, r := range allRoles {
			for _, ir := range svc.cfg.InitialIndexedRoles {
				if r.id == ir {
					out = append(out, r)
				}
			}
		}

		return
	}

	// Straight up limit
	// @todo add some counters to figure out which roles are most used from the start
	if svc.cfg.MaxIndexSize == -1 {
		return allRoles
	}

	if svc.cfg.MaxIndexSize == 0 {
		return nil
	}

	// @todo smarter way to figure out what/how many roles we want to load up
	return allRoles[:int(math.Min(float64(len(allRoles)), float64(svc.cfg.MaxIndexSize)))]
}
