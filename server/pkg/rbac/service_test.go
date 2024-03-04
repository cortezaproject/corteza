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

// Pre:
// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/pkg/rbac
// Benchmark_AccessCheck_role100_rule1000-12                  45502             26437 ns/op           10183 B/op        206 allocs/op
// Benchmark_AccessCheck_role100_rule10000-12                  9054            143348 ns/op           10195 B/op        206 allocs/op
// Benchmark_AccessCheck_role100_rule100000-12                  910           1399730 ns/op           10132 B/op        205 allocs/op
// Benchmark_AccessCheck_role100_rule1000000-12                  40          32568396 ns/op           11196 B/op        226 allocs/op
// Benchmark_AccessCheck_role100_rule10000000-12                  8         580995401 ns/op           14075 B/op        286 allocs/op
// Benchmark_AccessCheck_role1000_rule1000-12                 10000            115692 ns/op           78850 B/op       1216 allocs/op
// Benchmark_AccessCheck_role1000_rule10000-12                 4567            260073 ns/op           87800 B/op       1578 allocs/op
// Benchmark_AccessCheck_role1000_rule100000-12                 758           1521378 ns/op           87707 B/op       1593 allocs/op
// Benchmark_AccessCheck_role1000_rule1000000-12                 97          26178927 ns/op           79729 B/op       1447 allocs/op
// Benchmark_AccessCheck_role1000_rule10000000-12                 6         338761798 ns/op           87824 B/op       1502 allocs/op
// Benchmark_AccessCheck_role10000_rule1000-12                 1165           1113431 ns/op          875524 B/op      10173 allocs/op
// Benchmark_AccessCheck_role10000_rule10000-12                 972           1259840 ns/op          877246 B/op      11239 allocs/op
// Benchmark_AccessCheck_role10000_rule100000-12                404           2661588 ns/op          921948 B/op      14190 allocs/op
// Benchmark_AccessCheck_role10000_rule1000000-12               100          25640188 ns/op         1038714 B/op      15395 allocs/op
// Benchmark_AccessCheck_role10000_rule10000000-12                4         265117083 ns/op          882284 B/op      12844 allocs/op

// Post:
// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/pkg/rbac
// Benchmark_AccessCheck_role100_rule1000-12                 101308             11020 ns/op            5120 B/op         23 allocs/op
// Benchmark_AccessCheck_role100_rule10000-12                106794             10392 ns/op            5119 B/op         23 allocs/op
// Benchmark_AccessCheck_role100_rule100000-12               106137             10957 ns/op            5135 B/op         23 allocs/op
// Benchmark_AccessCheck_role100_rule1000000-12               94984             13411 ns/op            5110 B/op         23 allocs/op
// Benchmark_AccessCheck_role100_rule10000000-12              78258             13923 ns/op            5118 B/op         23 allocs/op
// Benchmark_AccessCheck_role1000_rule1000-12                 14739             79997 ns/op           52803 B/op         45 allocs/op
// Benchmark_AccessCheck_role1000_rule10000-12                10000            111318 ns/op           53129 B/op         45 allocs/op
// Benchmark_AccessCheck_role1000_rule100000-12               10000            119062 ns/op           53103 B/op         45 allocs/op
// Benchmark_AccessCheck_role1000_rule1000000-12              10000            127920 ns/op           53166 B/op         45 allocs/op
// Benchmark_AccessCheck_role1000_rule10000000-12              7801            406544 ns/op           53141 B/op         50 allocs/op
// Benchmark_AccessCheck_role10000_rule1000-12                 1609            744085 ns/op          508033 B/op        106 allocs/op
// Benchmark_AccessCheck_role10000_rule10000-12                1227            959672 ns/op          509599 B/op        106 allocs/op
// Benchmark_AccessCheck_role10000_rule100000-12                447           2711555 ns/op          502836 B/op        105 allocs/op
// Benchmark_AccessCheck_role10000_rule1000000-12               723           2202073 ns/op          528585 B/op        110 allocs/op
// Benchmark_AccessCheck_role10000_rule10000000-12              418           2587235 ns/op          496879 B/op        108 allocs/op
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
