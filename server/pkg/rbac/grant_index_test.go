package rbac

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// The index kept up to date grant by grant answers every lookup the way an
// index built from the current rules does: random allows, denies, changes and
// removals (Inherit), with and without a store to flush to
func TestGrantKeepsIndexInStep(t *testing.T) {
	var (
		ctx  = context.Background()
		rnd  = rand.New(rand.NewSource(42))
		ops  = []string{"read", "update"}
		ress = []string{
			"corteza::compose:record/*/*/*",
			"corteza::compose:record/1/*/*",
			"corteza::compose:record/1/2/*",
			"corteza::compose:record/1/2/3",
			"corteza::compose:record/1/2/4",
			"corteza::compose:module/1/2",
			"corteza::compose/",
		}
		access = []Access{Allow, Deny, Inherit}

		lookups = []string{
			"corteza::compose:record/1/2/3",
			"corteza::compose:record/1/2/4",
			"corteza::compose:record/1/2/5",
			"corteza::compose:record/9/9/9",
			"corteza::compose:module/1/2",
			"corteza::compose/",
			"corteza::compose",
		}
	)

	for _, store := range []rbacRulesStore{nil, nopRulesStore{}} {
		svc := NewService(zap.NewNop(), store)

		for i := 0; i < 400; i++ {
			rr := make([]*Rule, 1+rnd.Intn(3))
			for j := range rr {
				rr[j] = &Rule{
					RoleID:    uint64(1 + rnd.Intn(3)),
					Resource:  ress[rnd.Intn(len(ress))],
					Operation: ops[rnd.Intn(len(ops))],
					Access:    access[rnd.Intn(len(access))],
				}
			}
			require.NoError(t, svc.Grant(ctx, rr...))

			fresh := buildRuleIndex(svc.rules)
			for role := uint64(1); role <= 3; role++ {
				for _, op := range ops {
					for _, res := range lookups {
						msg := fmt.Sprintf("grant %d (store %T): role %d %s %s", i, store, role, op, res)
						require.Equal(t, ruleKeys(fresh.get(role, op, res)), ruleKeys(svc.indexed.get(role, op, res)), msg)

						want, got := fresh.matchingRule(role, op, res), svc.indexed.matchingRule(role, op, res)

						if want == nil {
							require.Nil(t, got, msg)
							continue
						}

						require.NotNil(t, got, msg)
						require.Equal(t, want.Access, got.Access, msg)
						require.Equal(t, want.Resource, got.Resource, msg)
					}
				}
			}
		}
	}
}

// ruleKeys lists the rules as comparable, sorted strings
func ruleKeys(rr []*Rule) (out []string) {
	for _, r := range rr {
		out = append(out, fmt.Sprintf("%d %s %s %s", r.RoleID, r.Operation, r.Resource, r.Access))
	}
	sort.Strings(out)
	return
}
