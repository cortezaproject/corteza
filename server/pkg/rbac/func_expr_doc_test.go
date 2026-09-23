package rbac

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"testing"

	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type rbacExprFunctionDoc struct {
	Name    string
	Example struct {
		Expr   string
		Scope  string
		Result string
	}
}

// loadRbacExprFunctionsDoc reads the functions expr_functions.yaml documents
// as registered by this package.
func loadRbacExprFunctionsDoc(t *testing.T) []rbacExprFunctionDoc {
	raw, err := os.ReadFile("../expr/expr_functions.yaml")
	require.NoError(t, err)

	var doc struct {
		Groups []struct {
			RegisteredIn string `yaml:"registeredIn"`
			Functions    []rbacExprFunctionDoc
		}
	}
	require.NoError(t, yaml.Unmarshal(raw, &doc))

	var out []rbacExprFunctionDoc
	for _, g := range doc.Groups {
		if g.RegisteredIn == "rbac" {
			out = append(out, g.Functions...)
		}
	}

	return out
}

// exampleOrgTree is the hierarchy the documented examples describe: user 1 in
// managers; user 2 in sales, under managers by a link named read; user 3 in
// support, under managers by an unnamed link.
func exampleOrgTree(t *testing.T) *orgTree {
	tree, err := mkOrgTree(
		ConvUserGroup(id.MustNumID(10), "managers", []id.ID{id.MustNumID(1)}, nil, nil),
		ConvUserGroup(id.MustNumID(20), "sales", []id.ID{id.MustNumID(2)}, nil,
			[]GroupNodePath{{SelfID: id.MustNumID(10), Name: "read"}}),
		ConvUserGroup(id.MustNumID(30), "support", []id.ID{id.MustNumID(3)}, nil,
			[]GroupNodePath{{SelfID: id.MustNumID(10)}}),
	)
	require.NoError(t, err)

	return tree
}

// TestRbacExprFunctionExamples evaluates the documented examples of this
// package's expression functions against the example hierarchy.
func TestRbacExprFunctionExamples(t *testing.T) {
	prevRBAC, prevChecker := gRBAC, isAboveChecker
	defer func() { gRBAC, isAboveChecker = prevRBAC, prevChecker }()

	gRBAC = &service{orgTree: exampleOrgTree(t)}
	isAboveChecker = func(owner id.ID, user id.ID, paths ...string) bool {
		return gRBAC.orgTree.IsAbove(owner, user, paths...)
	}

	expr.Init(AllFunctions, expr.AllFunctions)
	p := expr.NewParser()

	docs := loadRbacExprFunctionsDoc(t)
	require.NotEmpty(t, docs)

	for _, fn := range docs {
		fn := fn
		t.Run(fn.Name, func(t *testing.T) {
			ex := fn.Example
			require.NotEmptyf(t, ex.Expr, "%s has no example", fn.Name)
			require.NotEmptyf(t, ex.Result, "%s example has no result", fn.Name)

			var raw map[string]interface{}
			require.NoError(t, json.Unmarshal([]byte(ex.Scope), &raw))

			// IDs reach these functions as uint64, the way a contextual role's
			// scope carries them.
			scope := map[string]interface{}{"userID": uint64(raw["userID"].(float64))}
			if res, ok := raw["resource"].(map[string]interface{}); ok {
				scope["resource"] = map[string]interface{}{"ownedBy": uint64(res["ownedBy"].(float64))}
			}

			vars, err := expr.NewVars(scope)
			require.NoError(t, err)

			ev, err := p.Parse(ex.Expr)
			require.NoErrorf(t, err, "parse %q", ex.Expr)

			v, err := ev.Eval(context.Background(), vars)
			require.NoErrorf(t, err, "eval %q", ex.Expr)

			got, err := json.Marshal(expr.UntypedValue(v))
			require.NoError(t, err)
			require.JSONEqf(t, ex.Result, string(got), "%q", ex.Expr)
		})
	}
}

// TestRbacExprFunctionsDocumented holds expr_functions.yaml to the functions
// this package registers.
func TestRbacExprFunctionsDocumented(t *testing.T) {
	src, err := os.ReadFile("func_expr.go")
	require.NoError(t, err)

	var registered []string
	for _, m := range regexp.MustCompile(`gval\.Function\(\s*"([^"]+)"`).FindAllStringSubmatch(string(src), -1) {
		registered = append(registered, m[1])
	}
	require.NotEmpty(t, registered)

	var documented []string
	for _, fn := range loadRbacExprFunctionsDoc(t) {
		documented = append(documented, fn.Name)
	}

	sort.Strings(registered)
	sort.Strings(documented)
	require.Equal(t, registered, documented,
		"functions registered in func_expr.go and the expr_functions.yaml group registered in rbac differ")
}
