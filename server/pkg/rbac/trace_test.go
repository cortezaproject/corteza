package rbac

import (
	"testing"
)

func TestEvaluate(t *testing.T) {
	//var (
	//	rC1   = &Role{id: 1, kind: CommonRole}
	//	rAnon = &Role{id: 100, kind: AnonymousRole}
	//	rAuth = &Role{id: 200, kind: AuthenticatedRole}
	//	rCtx1 = &Role{id: 300, kind: ContextRole}
	//	rBps1 = &Role{id: 400, kind: BypassRole}
	//
	//	allRoles = []*Role{rC1, rAnon, rAuth, rCtx1, rBps1}
	//	_        = allRoles
	//
	//	res1 = "corteza::test:resource1"
	//
	//	tcc = []struct {
	//		name   string
	//		op     string
	//		res    string
	//		should Access
	//		roles  []*Role
	//		rules  RuleSet
	//	}{
	//		{
	//			name:  "base",
	//			op:    "READ",
	//			res:   res1,
	//			roles: []*Role{rAuth},
	//			rules: RuleSet{
	//				&Rule{RoleID: rAuth.id, Resource: res1, Operation: "READ", Access: Allow},
	//			},
	//
	//			should: Allow,
	//		},
	//	}
	//)
	//
	//for _, tc := range tcc {
	//	t.Run(tc.name, func(t *testing.T) {
	//		var (
	//			req = require.New(t)
	//
	//			partitionedRoles = partitionRoles(tc.roles...)
	//			indexedRules     = indexRules(tc.rules)
	//
	//			_access, _rule, _expl = evaluate(indexedRules, partitionedRoles, tc.op, tc.res, false)
	//		)
	//
	//		_ = _access
	//		_ = _rule
	//		_ = _expl
	//
	//		req.Equal(tc.should, _access)
	//	})
	//}
}
