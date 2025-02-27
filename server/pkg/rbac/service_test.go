package rbac

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	matchBenchCfg struct {
		rules RuleSet
		roles []*Role
		res   Resource
		op    string
	}

	stateCfg struct {
		resources       []string
		searchResponses []*Rule
	}
)

func TestNoopSvc(t *testing.T) {
	req := require.New(t)

	svc := NoopSvc(Allow, Config{})
	a, err := svc.Check(session{
		id:  1,
		rr:  []uint64{1},
		ctx: context.Background(),
	}, "read", NewResource("compose-record/1/2/3"))
	req.NoError(err)
	req.Equal(Allow, a)
}

func TestStatePrepping(t *testing.T) {
	ctx,
		req,
		svc, _ := prepState(
		t,
		stateCfg{
			resources: []string{
				"1:compose-record/1/2/3",
				"2:compose-record/1/2/3",
				"3:compose-record/1/2/3",
			},
			searchResponses: []*Rule{{
				RoleID:    1,
				Resource:  "compose-record/1/2/3",
				Operation: "read",
				Access:    Allow,
			}, {
				RoleID:    2,
				Resource:  "compose-record/1/2/3",
				Operation: "read",
				Access:    Deny,
			}},
		},
	)

	a, err := svc.Check(
		session{
			id:  1,
			rr:  []uint64{1},
			ctx: ctx,
		},
		"read",
		NewResource("compose-record/1/2/3"),
	)
	req.NoError(err)
	req.Equal(Allow, a)

	a, err = svc.Check(
		session{
			id:  1,
			rr:  []uint64{2},
			ctx: ctx,
		},
		"read",
		NewResource("compose-record/1/2/3"),
	)
	req.NoError(err)
	req.Equal(Deny, a)
}

func TestCheck_index(t *testing.T) {
	ctx,
		req,
		svc, _ := prepState(
		t,
		stateCfg{
			resources: []string{
				"1:compose-record/1/2/3",

				"2:compose-record/1/2/3",
				// "2:compose-record/1/2/*",

				"3:compose-record/1/2/3",
				// "3:compose-record/1/2/*",
				// "3:compose-record/1/*/*",

				"4:compose-record/1/2/3",
			},
			searchResponses: []*Rule{{
				RoleID:    1,
				Resource:  "compose-record/1/2/3",
				Operation: "read",
				Access:    Allow,
			},

				{
					RoleID:    2,
					Resource:  "compose-record/1/2/3",
					Operation: "read",
					Access:    Inherit,
				}, {
					RoleID:    2,
					Resource:  "compose-record/1/2/*",
					Operation: "read",
					Access:    Allow,
				},

				{
					RoleID:    3,
					Resource:  "compose-record/1/2/3",
					Operation: "read",
					Access:    Inherit,
				}, {
					RoleID:    3,
					Resource:  "compose-record/1/2/*",
					Operation: "read",
					Access:    Allow,
				}, {
					RoleID:    3,
					Resource:  "compose-record/1/*/*",
					Operation: "read",
					Access:    Deny,
				},

				{
					RoleID:    4,
					Resource:  "compose-record/1/2/3",
					Operation: "read",
					Access:    Deny,
				}},
		},
	)

	var a Access
	var err error

	t.Run("single rule, allow", func(t *testing.T) {
		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{1},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Allow, a)
	})

	t.Run("multi rule, inherit -> allow", func(t *testing.T) {
		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{2},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Allow, a)
	})

	t.Run("multi rule, allow because deny is lower", func(t *testing.T) {
		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{3},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Allow, a)
	})

	t.Run("single rule, deny", func(t *testing.T) {
		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{4},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Deny, a)
	})

	t.Run("multi role, deny higher up, deny", func(t *testing.T) {
		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{2, 4},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Deny, a)
	})

	t.Run("multi role, same level, deny over takes, deny", func(t *testing.T) {
		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{1, 4},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Deny, a)
	})
}

func TestIndexSize(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		_,
			req,
			svc, _ := prepState(
			t,
			stateCfg{
				resources: []string{
					"1:compose-record/1/2/3",
					"2:compose-record/1/2/3",
					"3:compose-record/1/2/3",
					"4:compose-record/1/2/3",
				},
				searchResponses: []*Rule{},
			},
		)

		req.Equal(0, svc.indexSize())
	})

	t.Run("all for same role", func(t *testing.T) {
		_,
			req,
			svc, _ := prepState(
			t,
			stateCfg{
				resources: []string{
					"1:compose-record/1/2/3",
					"2:compose-record/1/2/3",
					"3:compose-record/1/2/3",
					"4:compose-record/1/2/3",
				},
				searchResponses: []*Rule{{
					RoleID:    1,
					Resource:  "compose-record/1/2/3",
					Operation: "read",
					Access:    Allow,
				}},
			},
		)

		req.Equal(1, svc.indexSize())
	})

	t.Run("a bunch", func(t *testing.T) {
		_,
			req,
			svc, _ := prepState(
			t,
			stateCfg{
				resources: []string{
					"1:compose-record/1/2/3",
					"2:compose-record/1/2/3",
					"3:compose-record/1/2/3",
					"4:compose-record/1/2/3",
				},
				searchResponses: []*Rule{{
					RoleID:    1,
					Resource:  "compose-record/1/2/3",
					Operation: "read",
					Access:    Allow,
				},

					{
						RoleID:    2,
						Resource:  "compose-record/1/2/3",
						Operation: "read",
						Access:    Inherit,
					}, {
						RoleID:    2,
						Resource:  "compose-record/1/2/*",
						Operation: "read",
						Access:    Allow,
					},

					{
						RoleID:    3,
						Resource:  "compose-record/1/2/3",
						Operation: "read",
						Access:    Inherit,
					}, {
						RoleID:    3,
						Resource:  "compose-record/1/2/*",
						Operation: "read",
						Access:    Allow,
					}, {
						RoleID:    3,
						Resource:  "compose-record/1/*/*",
						Operation: "read",
						Access:    Deny,
					},
				},
			},
		)

		// Counting resources here!!
		req.Equal(3, svc.indexSize())
	})
}

func TestCheck_store(t *testing.T) {
	ctx,
		req,
		svc,
		store := prepState(
		t,
		stateCfg{
			resources: []string{
				// "1:compose-record/1/2/3",

				// "2:compose-record/1/2/3",
				// // "2:compose-record/1/2/*",

				// "3:compose-record/1/2/3",
				// // "3:compose-record/1/2/*",
				// // "3:compose-record/1/*/*",

				// "4:compose-record/1/2/3",
			},
		},
	)

	svc.roles = append(svc.roles,
		&Role{id: 1, handle: "1"},
		&Role{id: 2, handle: "2"},
		&Role{id: 3, handle: "3"},
		&Role{id: 4, handle: "4"},
	)

	var a Access
	var err error

	t.Run("single rule, allow", func(t *testing.T) {
		store.searchCount = 0
		store.searchResponses = [][]*Rule{{{
			RoleID:    1,
			Resource:  "compose-record/1/2/3",
			Operation: "read",
			Access:    Allow,
		}}}

		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{1},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Allow, a)
	})

	t.Run("multi rule, inherit -> allow", func(t *testing.T) {
		store.searchCount = 0
		store.searchResponses = [][]*Rule{{{
			RoleID:    2,
			Resource:  "compose-record/1/2/3",
			Operation: "read",
			Access:    Inherit,
		}, {
			RoleID:    2,
			Resource:  "compose-record/1/2/*",
			Operation: "read",
			Access:    Allow,
		}}}

		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{2},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Allow, a)
	})

	t.Run("multi rule, allow because deny is lower", func(t *testing.T) {
		store.searchCount = 0
		store.searchResponses = [][]*Rule{{{
			RoleID:    3,
			Resource:  "compose-record/1/2/3",
			Operation: "read",
			Access:    Inherit,
		}, {
			RoleID:    3,
			Resource:  "compose-record/1/2/*",
			Operation: "read",
			Access:    Allow,
		}, {
			RoleID:    3,
			Resource:  "compose-record/1/*/*",
			Operation: "read",
			Access:    Deny,
		}}}

		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{3},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Allow, a)
	})

	t.Run("single rule, deny", func(t *testing.T) {
		store.searchCount = 0
		store.searchResponses = [][]*Rule{{{
			RoleID:    4,
			Resource:  "compose-record/1/2/3",
			Operation: "read",
			Access:    Deny,
		}}}

		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{4},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Deny, a)
	})

	t.Run("multi role, deny higher up, deny", func(t *testing.T) {
		store.searchCount = 0
		store.searchResponses = [][]*Rule{{{
			RoleID:    2,
			Resource:  "compose-record/1/2/3",
			Operation: "read",
			Access:    Inherit,
		}, {
			RoleID:    2,
			Resource:  "compose-record/1/2/*",
			Operation: "read",
			Access:    Allow,
		}},

			{{
				RoleID:    4,
				Resource:  "compose-record/1/2/3",
				Operation: "read",
				Access:    Deny,
			}}}

		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{2, 4},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Deny, a)
	})

	t.Run("multi role, same level, deny over takes, deny", func(t *testing.T) {
		store.searchCount = 0
		store.searchResponses = [][]*Rule{{{
			RoleID:    1,
			Resource:  "compose-record/1/2/3",
			Operation: "read",
			Access:    Allow,
		}},

			{{
				RoleID:    4,
				Resource:  "compose-record/1/2/3",
				Operation: "read",
				Access:    Deny,
			}}}

		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{1, 4},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Deny, a)
	})
}

func TestCheck_mix(t *testing.T) {
	ctx,
		req,
		svc,
		store := prepState(
		t,
		stateCfg{
			resources: []string{
				// allow compose-record/1/2/3
				"1:compose-record/1/2/3",
				// deny compose-record/1/*/*
				"2:compose-record/1/2/3",
			},
			searchResponses: []*Rule{{
				RoleID:    1,
				Resource:  "compose-record/1/2/3",
				Operation: "read",
				Access:    Allow,
			},

				{
					RoleID:    2,
					Resource:  "compose-record/1/2/3",
					Operation: "read",
					Access:    Deny,
				},
			},
		},
	)

	store.searchResponse = nil

	svc.roles = append(svc.roles,
		&Role{id: 1, handle: "1"},
		&Role{id: 2, handle: "2"},
		&Role{id: 3, handle: "3"},
		&Role{id: 4, handle: "4"},

		&Role{id: 11, handle: "11"},
		&Role{id: 22, handle: "22"},
		&Role{id: 33, handle: "33"},
		&Role{id: 44, handle: "44"},
	)

	var a Access
	var err error

	t.Run("allow overtake inherit", func(t *testing.T) {
		store.searchCount = 0
		store.searchResponses = [][]*Rule{{{
			RoleID:    11,
			Resource:  "compose-record/1/2/3",
			Operation: "read",
			Access:    Inherit,
		}}}

		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{1, 11},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Allow, a)
	})

	t.Run("store deny overtake index allow", func(t *testing.T) {
		store.searchCount = 0
		store.searchResponses = [][]*Rule{{{
			RoleID:    11,
			Resource:  "compose-record/1/2/3",
			Operation: "read",
			Access:    Deny,
		}}}

		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{1, 11},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Deny, a)
	})

	t.Run("index deny overtakes db inherits", func(t *testing.T) {
		store.searchCount = 0
		store.searchResponses = [][]*Rule{{{
			RoleID:    11,
			Resource:  "compose-record/1/2/3",
			Operation: "read",
			Access:    Inherit,
		}, {
			RoleID:    11,
			Resource:  "compose-record/1/*/*",
			Operation: "read",
			Access:    Inherit,
		}, {
			RoleID:    11,
			Resource:  "compose-record/*/*/*",
			Operation: "read",
			Access:    Inherit,
		}}}

		a, err = svc.Check(
			session{
				id:  1,
				rr:  []uint64{2, 11},
				ctx: ctx,
			},
			"read",
			NewResource("compose-record/1/2/3"),
		)
		req.NoError(err)
		req.Equal(Deny, a)
	})
}

func prepState(t *testing.T, cfg stateCfg) (ctx context.Context, req *require.Assertions, svc *Service, store *mockRuleStore) {
	req = require.New(t)
	ctx = context.Background()

	store = &mockRuleStore{
		searchResponse: cfg.searchResponses,
		access:         Allow,
	}

	svc = &Service{
		RuleStorage: store,
		logger:      zap.NewNop(),
		roles:       rolesFromRes(cfg.resources...),
	}

	var err error
	svc.indexes, svc.indexMappings, err = svc.indexForResources(ctx, cfg.resources...)
	require.NoError(t, err)

	return
}

func rolesFromRes(resources ...string) (out []*Role) {
	for _, r := range resources {
		pp := strings.SplitN(r, ":", 2)
		out = append(out, CommonRole.Make(cast.ToUint64(pp[0]), pp[0]))
	}
	return
}

// goos: linux
// goarch: amd64
// pkg: github.com/cortezaproject/corteza/server/pkg/rbac
// cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
// Benchmark_AccessCheck_role5_rule500-12            378988              3026 ns/op             615 B/op         16 allocs/op
// Benchmark_AccessCheck_role5_rule1000-12           253071              4087 ns/op             615 B/op         16 allocs/op
// Benchmark_AccessCheck_role10_rule10000-12         237085              5429 ns/op            1026 B/op         29 allocs/op
// Benchmark_AccessCheck_role20_rule50000-12         128914              9344 ns/op            2335 B/op         71 allocs/op
// Benchmark_AccessCheck_role30_rule100000-12         79963             20670 ns/op            3371 B/op         85 allocs/op
// Benchmark_AccessCheck_role100_rule500000-12        16927             79106 ns/op           12796 B/op        391 allocs/op
// func benchmark_AccessCheck(b *testing.B, cfg matchBenchCfg) {
// 	svc := NewService(zap.NewNop(), nil)
// 	svc.UpdateRoles(cfg.roles...)
// 	svc.setRules(cfg.rules)

// 	ctx := context.Background()
// 	b.ResetTimer()

// 	for n := 0; n < b.N; n++ {
// 		svc.Can(session{
// 			id:  90001,
// 			rr:  yankRandRoles(cfg.roles),
// 			ctx: ctx,
// 		}, cfg.op, cfg.res)
// 	}
// }

// func Benchmark_AccessCheck_role100_rule1000(b *testing.B) {
// 	roles := 100
// 	rules := 1000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }

// func Benchmark_AccessCheck_role100_rule10000(b *testing.B) {
// 	roles := 100
// 	rules := 10000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }

// func Benchmark_AccessCheck_role100_rule100000(b *testing.B) {
// 	roles := 100
// 	rules := 100000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }

// func Benchmark_AccessCheck_role100_rule1000000(b *testing.B) {
// 	roles := 100
// 	rules := 1000000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }

// func Benchmark_AccessCheck_role100_rule10000000(b *testing.B) {
// 	roles := 100
// 	rules := 10000000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }

// func Benchmark_AccessCheck_role1000_rule1000(b *testing.B) {
// 	roles := 1000
// 	rules := 1000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }

// func Benchmark_AccessCheck_role1000_rule10000(b *testing.B) {
// 	roles := 1000
// 	rules := 10000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }

// func Benchmark_AccessCheck_role1000_rule100000(b *testing.B) {
// 	roles := 1000
// 	rules := 100000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }

// func Benchmark_AccessCheck_role1000_rule1000000(b *testing.B) {
// 	roles := 1000
// 	rules := 1000000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }

// func Benchmark_AccessCheck_role1000_rule10000000(b *testing.B) {
// 	roles := 1000
// 	rules := 10000000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }

// func Benchmark_AccessCheck_role10000_rule1000(b *testing.B) {
// 	roles := 10000
// 	rules := 1000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }
// func Benchmark_AccessCheck_role10000_rule10000(b *testing.B) {
// 	roles := 10000
// 	rules := 10000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }
// func Benchmark_AccessCheck_role10000_rule100000(b *testing.B) {
// 	roles := 10000
// 	rules := 100000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }
// func Benchmark_AccessCheck_role10000_rule1000000(b *testing.B) {
// 	roles := 10000
// 	rules := 1000000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }
// func Benchmark_AccessCheck_role10000_rule10000000(b *testing.B) {
// 	roles := 10000
// 	rules := 10000000
// 	benchmark_AccessCheck(b, matchBenchCfg{
// 		res:   makeResource(),
// 		op:    randomOperation(),
// 		rules: makeRuleSet(rules, roles),
// 		roles: makeRoleSet(roles),
// 	})
// }

// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/pkg/rbac
// cpu: Apple M3 Pro
// BenchmarkResourceBuild_100-12                180           6656197 ns/op          292020 B/op       4276 allocs/op
// BenchmarkResourceBuild_1000-12               126           9667167 ns/op         2846490 B/op      42264 allocs/op
// BenchmarkResourceBuild_10000-12               45          25771097 ns/op        32201567 B/op     422459 allocs/op
// BenchmarkResourceBuild_100000-12               2         693479625 ns/op        960572552 B/op  11028408 allocs/op
func BenchmarkResourceBuild_100(b *testing.B)    { benchmarkResourceBuild(b, 100) }
func BenchmarkResourceBuild_1000(b *testing.B)   { benchmarkResourceBuild(b, 1000) }
func BenchmarkResourceBuild_10000(b *testing.B)  { benchmarkResourceBuild(b, 10000) }
func BenchmarkResourceBuild_100000(b *testing.B) { benchmarkResourceBuild(b, 100000) }

func benchmarkResourceBuild(b *testing.B, n int) {
	ctx := context.Background()
	svc := &Service{
		RuleStorage: &mockRuleStore{},
	}

	var res []string
	roleID := uint64(0)
	for i := 0; i < n; i++ {
		if i%2000 == 0 {
			roleID++
		}

		res = append(res, fmt.Sprintf("%d:compose::record/%d/%d/%d", roleID, i, i, i))
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		svc.indexForResources(ctx, res...)
	}
}

func yankRandRoles(base []*Role) (out []uint64) {
	count := rand.Intn(len(base))
	if count <= 0 {
		count = int(math.Ceil(float64(len(base)) / 2))
	}

	out = make([]uint64, count)
	for i := 0; i < count; i++ {
		out[i] = base[i].id
	}

	return
}

func makeRoleSet(count int) (out []*Role) {
	for i := 0; i < count; i++ {
		out = append(out, makeRole(uint64(1000+i), fmt.Sprintf("rl_%d", 1000+i)))
	}

	return
}

func makeRole(id uint64, handle string) *Role {
	rx := rand.Float64()

	if rx <= 1 {
		return CommonRole.Make(id, handle)
	}

	return makeContextualRole(id, handle)
}

func makeContextualRole(id uint64, handle string) *Role {
	rx := rand.Float64()

	if rx < 0.7 {
		return makeContextualRolePassing(id, handle)
	}
	return makeContextualRoleFailing(id, handle)
}

func makeContextualRolePassing(id uint64, handle string) *Role {
	p := expr.NewParser()
	eval, err := p.Parse("true == true && true == true && 1 <= 1")
	if err != nil {
		panic(err)
	}

	check := func(scope map[string]interface{}) bool {
		vars, err := expr.NewVars(scope)
		if err != nil {
			return false
		}

		ctx := context.Background()
		test, err := eval.Test(ctx, vars)
		if err != nil {
			return false
		}

		return test
	}

	return MakeContextRole(id, handle, check, "corteza::compose:record")
}

func makeContextualRoleFailing(id uint64, handle string) *Role {
	p := expr.NewParser()
	eval, err := p.Parse("false == false || false == false || 1 > 1")
	if err != nil {
		panic(err)
	}

	check := func(scope map[string]interface{}) bool {
		vars, err := expr.NewVars(scope)
		if err != nil {
			return false
		}

		ctx := context.Background()
		test, err := eval.Test(ctx, vars)
		if err != nil {
			return false
		}

		return test
	}

	return MakeContextRole(id, handle, check, "corteza::compose:record")
}

func makeResource() (out Resource) {
	return resource(randomResource())
}

func randomAccess() (out Access) {
	x := rand.Float64()
	if x < 0.7 {
		return Inherit
	}
	return Inherit
}

func randomOperation() (out string) {
	ops := []string{"read", "write", "delete"}
	return ops[rand.Intn(len(ops))]
}

func randomResource() (out string) {
	return fmt.Sprintf("%s:%s/%s/%s", RandStringRunes(1), RandStringRunes(1), RandStringRunes(1), RandStringRunes(1))
}

// var letterRunes = []rune("abcdefghijklmnoprst")
var letterRunes = []rune("abcdefghijklmnoprstuvzxyABCDEFGHIJKLMNOPRSTXY")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
