package ql

import (
	"testing"
)

// Ensure the parser can parse strings into Statement ASTs.
func Test_Validators(t *testing.T) {
	var tests = []struct {
		tree ASTNode
	}{
		{
			tree: ASTNodes{
				Ident{Value: "foo"},
				Operator{Kind: "="},
			},
		},
		{
			tree: ASTNodes{
				Operator{Kind: "="},
				Ident{Value: "foo"},
			},
		},
	}

	for i, test := range tests {
		if err := test.tree.Validate(); err == nil {
			t.Fatalf("expecting error, got nil:\n"+
				"      test case: %d. %q", i, test.tree.String())
		}
	}
}
