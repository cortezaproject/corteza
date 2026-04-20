package expr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Example_simpleExpression() {
	eval(`40 + 2`, nil)
	// output:
	// 42
}

func TestParser(t *testing.T) {
	var (
		req    = require.New(t)
		ctx    = context.Background()
		p      = Parser()
		e, err = p.NewEvaluable("0 == 0")

		result bool
	)

	req.NoError(err)

	result, err = e.EvalBool(ctx, nil)
	req.NoError(err)
	req.True(result)
}

func TestGvalParser(t *testing.T) {
	var (
		req     = require.New(t)
		ctx     = context.Background()
		p       = NewGvalParser()
		vv, err = NewVars(map[string]interface{}{
			"vars":  &Vars{},
			"key":   "foo",
			"value": Must(NewString("foo")),
		})
		result interface{}
	)
	req.NoError(err)

	pp, err := p.Parse("toJSON(set(vars, key, value))")
	req.NoError(err)

	result, err = pp.Eval(ctx, vv)
	req.NoError(err)
	req.Equal("{\"foo\":{\"@value\":\"foo\",\"@type\":\"String\"}}", result)
}

// goos: darwin
// goarch: arm64
// pkg: github.com/crusttech/human/server/pkg/expr
// cpu: Apple M3 Pro
// BenchmarkParsing-12    	  295779	      4052 ns/op	    4952 B/op	     133 allocs/op
func BenchmarkParsing(b *testing.B) {
	expr := `(1 + length(trim(" asdf asdf asdf ")) * 71 + min(3, 7777777)) + (1 + length(trim(" asdf asdf asdf ")) * 71 + min(3, 7777777))`
	p := NewGvalParser()

	for n := 0; n < b.N; n++ {
		e, err := p.Parse(expr)
		if err != nil {
			panic(err)
		}
		_ = e
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/crusttech/human/server/pkg/expr
// cpu: Apple M3 Pro
// BenchmarkEval_short-12            886696              1344 ns/op             544 B/op         14 allocs/op
// BenchmarkEval_med-12              299616              4020 ns/op            1280 B/op         36 allocs/op
// BenchmarkEval_long-12             146112              7878 ns/op            2488 B/op         71 allocs/op
func BenchmarkEval_short(b *testing.B) {
	expr := `min(3, 7777777)`
	benchmarkEval(b, expr)
}

func BenchmarkEval_med(b *testing.B) {
	expr := `1 + length(trim(" asdf asdf asdf ")) * 71 + min(3, 7777777)`
	benchmarkEval(b, expr)
}

func BenchmarkEval_long(b *testing.B) {
	expr := `(1 + length(trim(" asdf asdf asdf ")) * 71 + min(3, 7777777)) + (1 + length(trim(" asdf asdf asdf ")) * 71 + min(3, 7777777))`
	benchmarkEval(b, expr)
}

func benchmarkEval(b *testing.B, expr string) {
	e, err := NewGvalParser().Parse(expr)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	for n := 0; n < b.N; n++ {
		_, err = e.Eval(ctx, EmptyVars())
		// _ = err
		if err != nil {
			panic(err)
		}
	}
}
