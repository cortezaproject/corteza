package agentic

import (
	"errors"
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/require"
)

func TestWithValueIssues(t *testing.T) {
	base := errors.New("invalid record value input")

	t.Run("no issues leaves the error alone", func(t *testing.T) {
		require.Equal(t, base, withValueIssues(base, nil))
		require.Equal(t, base, withValueIssues(base, &cmpTypes.RecordValueErrorSet{}))
	})

	t.Run("names the field and the rule", func(t *testing.T) {
		dd := &cmpTypes.RecordValueErrorSet{Set: []cmpTypes.RecordValueError{
			{Kind: "empty", Message: "This field is required", Meta: map[string]interface{}{"field": "title"}},
		}}

		err := withValueIssues(base, dd)
		require.ErrorIs(t, err, base, "the original error must stay in the chain")
		require.Contains(t, err.Error(), "title")
		require.Contains(t, err.Error(), "This field is required")
	})

	t.Run("joins several, and copes with a missing field name", func(t *testing.T) {
		dd := &cmpTypes.RecordValueErrorSet{Set: []cmpTypes.RecordValueError{
			{Message: "This field is required", Meta: map[string]interface{}{"field": "title"}},
			{Message: "duplicate value", Meta: map[string]interface{}{}},
		}}

		require.Equal(t,
			"invalid record value input (title: This field is required; duplicate value)",
			withValueIssues(base, dd).Error(),
		)
	})
}
