package rbac

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"go.uber.org/zap"
)

type (
	matchBenchCfg struct {
		rules RuleSet
		roles []*Role
		res   Resource
		op    string
	}
)

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
func benchmark_AccessCheck(b *testing.B, cfg matchBenchCfg) {
	svc := NewService(zap.NewNop(), nil)
	svc.UpdateRoles(cfg.roles...)
	svc.setRules(cfg.rules)

	ctx := context.Background()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		svc.Can(session{
			id:  90001,
			rr:  yankRandRoles(cfg.roles),
			ctx: ctx,
		}, cfg.op, cfg.res)
	}
}

func Benchmark_AccessCheck_role100_rule1000(b *testing.B) {
	roles := 100
	rules := 1000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}

func Benchmark_AccessCheck_role100_rule10000(b *testing.B) {
	roles := 100
	rules := 10000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}

func Benchmark_AccessCheck_role100_rule100000(b *testing.B) {
	roles := 100
	rules := 100000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}

func Benchmark_AccessCheck_role100_rule1000000(b *testing.B) {
	roles := 100
	rules := 1000000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}

func Benchmark_AccessCheck_role100_rule10000000(b *testing.B) {
	roles := 100
	rules := 10000000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}

func Benchmark_AccessCheck_role1000_rule1000(b *testing.B) {
	roles := 1000
	rules := 1000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}

func Benchmark_AccessCheck_role1000_rule10000(b *testing.B) {
	roles := 1000
	rules := 10000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}

func Benchmark_AccessCheck_role1000_rule100000(b *testing.B) {
	roles := 1000
	rules := 100000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}

func Benchmark_AccessCheck_role1000_rule1000000(b *testing.B) {
	roles := 1000
	rules := 1000000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}

func Benchmark_AccessCheck_role1000_rule10000000(b *testing.B) {
	roles := 1000
	rules := 10000000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}

func Benchmark_AccessCheck_role10000_rule1000(b *testing.B) {
	roles := 10000
	rules := 1000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}
func Benchmark_AccessCheck_role10000_rule10000(b *testing.B) {
	roles := 10000
	rules := 10000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}
func Benchmark_AccessCheck_role10000_rule100000(b *testing.B) {
	roles := 10000
	rules := 100000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}
func Benchmark_AccessCheck_role10000_rule1000000(b *testing.B) {
	roles := 10000
	rules := 1000000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
}
func Benchmark_AccessCheck_role10000_rule10000000(b *testing.B) {
	roles := 10000
	rules := 10000000
	benchmark_AccessCheck(b, matchBenchCfg{
		res:   makeResource(),
		op:    randomOperation(),
		rules: makeRuleSet(rules, roles),
		roles: makeRoleSet(roles),
	})
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

func makeRuleSet(count int, roleCount int) (out RuleSet) {
	for i := 0; i < count; i++ {
		out = append(out, &Rule{
			RoleID:    uint64(1000 + int(rand.Intn(roleCount))),
			Resource:  randomResource(),
			Operation: randomOperation(),
			Access:    randomAccess(),
		})
	}

	return
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
