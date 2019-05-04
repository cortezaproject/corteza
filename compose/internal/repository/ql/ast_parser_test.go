// +build unit

package ql

import (
	"fmt"
	"reflect"
	"testing"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestAstParser_Parser(t *testing.T) {
	var tests = []struct {
		in     string
		tree   ASTNode
		err    error
		sql    string
		args   []interface{}
		parser func(s string) (ASTNode, error)
	}{
		{
			in: `log( arg1 ), arg2 / 100`,
			tree: ASTSet{
				Function{
					Name: "log",
					Arguments: ASTSet{
						Ident{Value: "arg1"},
					},
				},

				ASTNodes{
					Ident{Value: "arg2"},
					Operator{Kind: "/"},
					Number{Value: "100"},
				},
			},
		},
		{
			in: `log( arg1 * 5 ), arg2 / 100 + 10`,
			tree: ASTSet{
				Function{
					Name: "log",
					Arguments: ASTSet{
						ASTNodes{
							Ident{Value: "arg1"},
							Operator{Kind: "*"},
							Number{Value: "5"},
						},
					},
				},

				ASTNodes{
					Ident{Value: "arg2"},
					Operator{Kind: "/"},
					Number{Value: "100"},
					Operator{Kind: "+"},
					Number{Value: "10"},
				},
			},
		},
		{
			in: `date_format(created_at, '%Y')`,
			tree: ASTSet{
				Function{
					Name: "date_format",
					Arguments: ASTSet{
						Ident{Value: "created_at"},
						String{Value: "%Y"},
					},
				},
			},
			sql:  `date_format(created_at, ?)`,
			args: []interface{}{"%Y"},
		},
		{
			parser: NewParser().ParseExpression,
			in:     `func(arg1, arg2)`,
			tree: Function{
				Name: "func",
				Arguments: ASTSet{
					Ident{Value: "arg1"},
					Ident{Value: "arg2"},
				},
			},
		},
		{
			parser: NewParser().ParseExpression,
			in:     `year(created_at) != 2010`,
			tree: ASTNodes{
				Function{
					Name: "year",
					Arguments: ASTSet{
						Ident{Value: "created_at"},
					},
				},
				Operator{Kind: "!="},
				Number{Value: "2010"},
			},
		},
		{
			parser: NewParser().ParseExpression,
			in:     `year(created_at) != 2010 AND month(created_at) = 6`,
			tree: ASTNodes{
				Function{
					Name: "year",
					Arguments: ASTSet{
						Ident{Value: "created_at"},
					},
				},
				Operator{Kind: "!="},
				Number{Value: "2010"},
				Operator{"AND"},
				Function{
					Name: "month",
					Arguments: ASTSet{
						Ident{Value: "created_at"},
					},
				},
				Operator{Kind: "="},
				Number{Value: "6"},
			},
		},
		{
			parser: NewParser().ParseExpression,
			in:     `year(created_at) = year(now()) - 1`,
			tree: ASTNodes{
				Function{Name: "year", Arguments: ASTSet{Ident{Value: "created_at"}}},
				Operator{Kind: "="},
				Function{Name: "year", Arguments: ASTSet{Function{Name: "now"}}},
				Operator{Kind: "-"},
				Number{Value: "1"},
			},
		},
		{
			parser: NewParser().ParseExpression,
			in:     `NOW() > DATE_SUB(col, INTERVAL 31 DAY)`,
			tree: ASTNodes{
				Function{Name: "NOW"},
				Operator{Kind: ">"},
				Function{Name: "DATE_SUB", Arguments: ASTSet{Ident{Value: "col"}, Interval{Value: "31", Unit: "DAY"}}},
			},
		},
		{
			parser: NewParser().ParseExpression,
			in:     `foo LIKE 'bar%'`,
			tree: ASTNodes{
				Ident{Value: "foo"},
				Operator{Kind: "LIKE"},
				String{Value: "bar%"},
			},
		},
		{
			parser: NewParser().ParseExpression,
			in:     `foo NOT LIKE 'bar%'`,
			tree: ASTNodes{
				Ident{Value: "foo"},
				Operator{Kind: "NOT LIKE"},
				String{Value: "bar%"},
			},
		},
	}

	for i, test := range tests {
		if test.parser == nil {
			test.parser = NewParser().ParseSet
		}

		if tree, err := test.parser(test.in); err != test.err {
			t.Fatalf("error mismatch:\n"+
				"test case: %d. %s\n"+
				" expected: %v\n"+
				"      got: %v\n\n", i, test.in, test.err, err)
		} else if test.err == nil && !reflect.DeepEqual(test.tree, tree) {
			t.Errorf("tree does not match:\n"+
				"test case: %d. %s\n"+
				" expected: %#v\n"+
				"      got: %#v\n\n", i, test.in, test.tree, tree)
		} else if sql, args, err := tree.ToSql(); err != nil {
			t.Fatal(err)
		} else if test.sql != "" && sql != test.sql {
			t.Errorf("sql does not match:\n"+
				"test case: %d. %s\n"+
				" expected: %#v\n"+
				"      got: %#v\n\n", i, test.in, test.sql, sql)
		} else if test.args != nil && !reflect.DeepEqual(test.args, args) {
			t.Errorf("args does not match:\n"+
				"test case: %d. %s\n"+
				" expected: %#v\n"+
				"      got: %#v\n\n", i, test.in, test.args, args)
		}
	}
}

func TestAstParser_ColumnParser(t *testing.T) {
	var tests = []struct {
		in   string
		cols Columns
		err  error
	}{
		{
			in: `a AS b`,
			cols: Columns{
				Column{
					Expr:  ASTNodes{Ident{Value: "a"}},
					Alias: "b",
				},
			},
		},
		{
			in: `sum(value1) as sumValue1, min(value2)`,
			cols: Columns{
				Column{
					Expr: ASTNodes{Function{
						Name: "sum",
						Arguments: ASTSet{
							Ident{Value: "value1"},
						},
					}},
					Alias: "sumValue1",
				},
				Column{
					Expr: ASTNodes{Function{
						Name: "min",
						Arguments: ASTSet{
							Ident{Value: "value2"},
						},
					}},
				},
			},
		},
		{
			in: `a DESC`,
			cols: Columns{
				Column{
					Expr: ASTNodes{Ident{Value: "a"}, Keyword{Keyword: "DESC"}},
				},
			},
		},
		{
			in: `a ASC`,
			cols: Columns{
				Column{
					Expr: ASTNodes{Ident{Value: "a"}, Keyword{Keyword: "ASC"}},
				},
			},
		},
		{
			in: `DATE_FORMAT(some_date, '%Y-%m-01')`,
			cols: Columns{
				Column{
					Expr: ASTNodes{
						Function{
							Name: "DATE_FORMAT",
							Arguments: ASTSet{
								Ident{Value: "some_date"},
								String{Value: "%Y-%m-01"},
							},
						},
					},
				},
			},
		},
	}

	p := NewParser()
	for i, test := range tests {
		if cols, err := p.ParseColumns(test.in); err != test.err {
			t.Fatalf("%d. %s: error mismatch:\n  expected: %v\n        got: %v\n\n", i, test.in, test.err, err)
		} else if test.err == nil && !reflect.DeepEqual(test.cols, cols) {
			t.Errorf("%d. %s\n\ncols does not match:\n\nexpected: %#v\n     got: %#v\n\n", i, test.in, test.cols, cols)
		}
	}
}

func TestAstParser_IdentModifier(t *testing.T) {
	var tests = []struct {
		in  string
		out string
		err error
	}{
		{
			in:  `foo`,
			out: `__wrap_foo_wrap__`,
		},
	}

	p := NewParser()
	for i, test := range tests {
		p.OnIdent = func(ident Ident) (Ident, error) {
			ident.Value = fmt.Sprintf("__wrap_%s_wrap__", ident.Value)
			return ident, nil
		}

		if tree, err := p.ParseExpression(test.in); err != test.err {
			t.Fatalf("%d. error mismatch:\n  expected: %v\n        got: %v\n\n", i, test.err, err)
		} else if test.err == nil && test.out != tree.String() {
			t.Errorf("%d. tree does not match:\n\n expected: %#v\n      got: %#v\n\n", i, test.in, test.out)
		}
	}
}
