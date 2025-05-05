package rbac

import (
	"context"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/spf13/cast"
)

type (
	mockRuleStore struct {
		searchResponse []*Rule

		searchResponses [][]*Rule
		searchCount     int

		searches []RuleFilter
		access   Access
	}
)

func (svc *mockRuleStore) SearchRbacRules(ctx context.Context, f RuleFilter) (out RuleSet, _ RuleFilter, err error) {
	svc.searches = append(svc.searches, f)

	if svc.searchResponses != nil {
		out = svc.searchResponses[svc.searchCount]
		svc.searchCount++
		return
	}

	if svc.searchResponse != nil {
		return svc.searchResponse, f, nil
	}

	if f.RawFilter != "" {
		// (rel_role=%d and (%s))
		combos := strings.Split(f.RawFilter, ") or (")

		for _, c := range combos {
			roleRules := strings.Split(c, " and (")

			role := strings.Split(roleRules[0], "=")[1]

			for _, rule := range strings.Split(roleRules[1], " or ") {
				rs := strings.Split(rule, "resource=")[1]

				out = append(out, &Rule{
					RoleID:    cast.ToUint64(role),
					Resource:  rs[1 : len(rs)-1],
					Operation: "read",
					Access:    svc.access,
				})
			}
		}
	} else {
		for _, r := range f.Resource {
			out = append(out, &Rule{
				RoleID:    f.RoleID,
				Resource:  r,
				Operation: f.Operation,
				Access:    svc.access,
			})
		}
	}

	// Give a slight delay to better simulate what we'd expect to see in real-world
	time.Sleep(time.Millisecond * 5)
	return
}

func (svc mockRuleStore) UpsertRbacRule(ctx context.Context, rr ...*Rule) (err error) {
	return
}

func (svc mockRuleStore) DeleteRbacRule(ctx context.Context, rr ...*Rule) (err error) {
	return
}

func (svc mockRuleStore) TruncateRbacRules(ctx context.Context) (err error) {
	return
}

func (svc mockRuleStore) SearchRoles(ctx context.Context, f types.RoleFilter) (out types.RoleSet, _ types.RoleFilter, err error) {
	return
}
