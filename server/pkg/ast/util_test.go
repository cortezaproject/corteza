package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeAnd(t *testing.T) {
	a := &ASTNode{Ref: "eq", Args: ASTNodeSet{{Ref: "x"}, {Ref: "1"}}}
	b := &ASTNode{Ref: "eq", Args: ASTNodeSet{{Ref: "y"}, {Ref: "2"}}}

	t.Run("both non-nil", func(t *testing.T) {
		result := MergeAnd(a, b)
		assert.Equal(t, "and", result.Ref)
		assert.Len(t, result.Args, 2)
		assert.Equal(t, "group", result.Args[0].Ref)
		assert.Equal(t, "group", result.Args[1].Ref)
	})

	t.Run("a nil", func(t *testing.T) {
		result := MergeAnd(nil, b)
		assert.Equal(t, b, result)
	})

	t.Run("b nil", func(t *testing.T) {
		result := MergeAnd(a, nil)
		assert.Equal(t, a, result)
	})

	t.Run("both nil", func(t *testing.T) {
		result := MergeAnd(nil, nil)
		assert.Nil(t, result)
	})
}

func TestMergeOr(t *testing.T) {
	a := &ASTNode{Ref: "eq", Args: ASTNodeSet{{Ref: "x"}, {Ref: "1"}}}
	b := &ASTNode{Ref: "eq", Args: ASTNodeSet{{Ref: "y"}, {Ref: "2"}}}

	t.Run("both non-nil", func(t *testing.T) {
		result := MergeOr(a, b)
		assert.Equal(t, "or", result.Ref)
		assert.Len(t, result.Args, 2)
		assert.Equal(t, "group", result.Args[0].Ref)
		assert.Equal(t, "group", result.Args[1].Ref)
	})

	t.Run("a nil", func(t *testing.T) {
		result := MergeOr(nil, b)
		assert.Equal(t, b, result)
	})

	t.Run("b nil", func(t *testing.T) {
		result := MergeOr(a, nil)
		assert.Equal(t, a, result)
	})

	t.Run("both nil", func(t *testing.T) {
		result := MergeOr(nil, nil)
		assert.Nil(t, result)
	})
}
