package values

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/stretchr/testify/require"
)

func TestExpressions(t *testing.T) {
	var (
		ctx = context.Background()

		// pairs = [<fieldname>, <value expression>, ...]
		makeModule = func(pairs ...string) *types.Module {
			var (
				m = &types.Module{}
			)

			for i := 0; i < len(pairs); i += 3 {
				f := &types.ModuleField{Name: pairs[i], Kind: pairs[i+1], Options: map[string]interface{}{}}
				f.Expressions.ValueExpr = pairs[i+2]
				m.Fields = append(m.Fields, f)
			}

			return m
		}
	)

	t.Run("empty", func(t *testing.T) {
		var (
			req = require.New(t)
			m   = &types.Module{}
			r   = &types.Record{}
			rve = &types.RecordValueErrorSet{}
		)

		Expression(ctx, m, r, nil, rve)
		req.True(rve.IsValid())
	})

	t.Run("string", func(t *testing.T) {
		var (
			req = require.New(t)
			m   = makeModule("f1", "String", `"abc"`)
			r   = &types.Record{}
			rve = &types.RecordValueErrorSet{}
		)

		Expression(ctx, m, r, nil, rve)
		req.Truef(rve.IsValid(), "%v", rve.Set)
		req.Equal("abc", r.Values.Get("f1", 0).Value)
	})

	t.Run("use fields", func(t *testing.T) {
		var (
			req = require.New(t)
			m   = makeModule("f1", "String", `fname + " " + lname`)
			r   = &types.Record{}
			rve = &types.RecordValueErrorSet{}
		)

		m.Fields = append(m.Fields, &types.ModuleField{Name: "fname", Kind: "String"})
		m.Fields = append(m.Fields, &types.ModuleField{Name: "lname", Kind: "String"})
		r.Values = r.Values.Replace("fname", "Cor")
		r.Values = r.Values.Replace("lname", "Teza")

		Expression(ctx, m, r, nil, rve)
		req.Truef(rve.IsValid(), "%v", rve.Set)
		req.Equal("Cor Teza", r.Values.Get("f1", 0).Value)
	})

	t.Run("math", func(t *testing.T) {
		var (
			req = require.New(t)
			m   = makeModule("f1", "Number", `n1 * n2 * 1.511`)
			r   = &types.Record{}
			rve = &types.RecordValueErrorSet{}
		)

		m.Fields.FindByName("f1").Options.SetPrecision(1)

		m.Fields = append(m.Fields, &types.ModuleField{Name: "n1", Kind: "Number"})
		m.Fields = append(m.Fields, &types.ModuleField{Name: "n2", Kind: "Number"})
		r.Values = r.Values.Replace("n1", "10")
		r.Values = r.Values.Replace("n2", "20")

		Expression(ctx, m, r, nil, rve)
		req.Truef(rve.IsValid(), "%v", rve.Set)
		req.Equal("302.2", r.Values.Get("f1", 0).Value)
	})

	t.Run("booleans", func(t *testing.T) {
		var (
			m   = makeModule("f1", "string", `b1 ? "yes" : "no"`)
			r   = &types.Record{}
			rve = &types.RecordValueErrorSet{}
		)

		m.Fields = append(m.Fields, &types.ModuleField{Name: "b1", Kind: "Bool"})

		t.Run("true", func(t *testing.T) {
			var req = require.New(t)
			r.Values = r.Values.Replace("b1", "true")

			Expression(ctx, m, r, nil, rve)
			req.Truef(rve.IsValid(), "%v", rve.Set)
			req.Equal("yes", r.Values.Get("f1", 0).Value)
		})

		t.Run("false", func(t *testing.T) {
			var req = require.New(t)
			r.Values = r.Values.Replace("b1", "false")

			Expression(ctx, m, r, nil, rve)
			req.Truef(rve.IsValid(), "%v", rve.Set)
			req.Equal("no", r.Values.Get("f1", 0).Value)
		})
	})

	t.Run("old vs new", func(t *testing.T) {
		var (
			req = require.New(t)
			m   = makeModule("f1", "String", `new.values.test == old.values.test ? "same":"different"`)
			new = &types.Record{}
			old = &types.Record{}
			rve = &types.RecordValueErrorSet{}
		)

		m.Fields = append(m.Fields, &types.ModuleField{Name: "test", Kind: "String"})

		new.Values = new.Values.Replace("test", "a")
		old.Values = old.Values.Replace("test", "a")

		Expression(ctx, m, new, old, rve)
		req.Truef(rve.IsValid(), "%v", rve.Set)
		req.Equal("same", new.Values.Get("f1", 0).Value)
	})

	t.Run("multi value expressions", func(t *testing.T) {
		var (
			req = require.New(t)
			m   = makeModule("f1", "String", `[t1,t2,t3]`)
			new = &types.Record{}
			old = &types.Record{}
			rve = &types.RecordValueErrorSet{}
		)

		m.Fields.FindByName("f1").Multi = true

		m.Fields = append(m.Fields, &types.ModuleField{Name: "t1", Kind: "String"})
		m.Fields = append(m.Fields, &types.ModuleField{Name: "t2", Kind: "String"})
		m.Fields = append(m.Fields, &types.ModuleField{Name: "t3", Kind: "String"})

		new.Values = new.Values.Replace("t1", "a")
		new.Values = new.Values.Replace("t2", "b")
		new.Values = new.Values.Replace("t3", "c")

		Expression(ctx, m, new, old, rve)
		req.Truef(rve.IsValid(), "%v", rve.Set)
		req.Equal("a", new.Values.Get("f1", 0).Value)
		req.Equal("b", new.Values.Get("f1", 1).Value)
		req.Equal("c", new.Values.Get("f1", 2).Value)
	})

	t.Run("omit record value when expr returns null (nil)", func(t *testing.T) {
		var (
			req = require.New(t)
			m   = makeModule("f1", "String", `null`)
			r   = &types.Record{}
			rve = &types.RecordValueErrorSet{}
		)

		Expression(ctx, m, r, nil, rve)
		req.True(rve.IsValid())
		req.Empty(r.Values)
	})
}
