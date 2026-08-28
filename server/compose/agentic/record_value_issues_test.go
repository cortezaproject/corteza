package agentic

import (
	"errors"
	"strings"
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	hErrors "github.com/crusttech/human/server/pkg/errors"
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

	// The validation path never returns the set; it wraps it into the error, so
	// a fix that only reads the return value leaves a missing required field
	// reported as a bare "invalid record value input".
	t.Run("finds the set wrapped in the error", func(t *testing.T) {
		dd := &cmpTypes.RecordValueErrorSet{Set: []cmpTypes.RecordValueError{
			{Kind: "empty", Message: "This field is required", Meta: map[string]interface{}{"field": "title"}},
		}}
		wrapped := hErrors.New(hErrors.KindInternal, "invalid record value input").Wrap(dd)

		err := withValueIssues(wrapped, nil)
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

func TestFillValuePlaceholders(t *testing.T) {
	t.Run("the value comes out of the meta beside it", func(t *testing.T) {
		require.Equal(t,
			`The value "taken@example.com" already exists in another record`,
			fillValuePlaceholders(
				`The value "{{value}}" already exists in another record`,
				map[string]any{"value": "taken@example.com"},
			),
		)
	})

	t.Run("a slot with nothing to fill it is left as written", func(t *testing.T) {
		require.Equal(t,
			`The value "{{value}}" already exists`,
			fillValuePlaceholders(`The value "{{value}}" already exists`, map[string]any{"field": "email"}),
		)
	})

	t.Run("whitespace inside the slot is still a slot", func(t *testing.T) {
		require.Equal(t, "got 7", fillValuePlaceholders("got {{ n }}", map[string]any{"n": 7}))
	})

	t.Run("a message with no slots is untouched", func(t *testing.T) {
		require.Equal(t, "This field is required", fillValuePlaceholders("This field is required", map[string]any{"value": "x"}))
	})
}

func TestWithValueIssuesFillsThePlaceholder(t *testing.T) {
	base := errors.New("invalid record value input")
	dd := &cmpTypes.RecordValueErrorSet{Set: []cmpTypes.RecordValueError{{
		Kind:    "duplicateValue",
		Message: `The value "{{value}}" already exists in another record`,
		Meta:    map[string]interface{}{"field": "email", "value": "a@example.com"},
	}}}

	err := withValueIssues(base, dd)
	require.Contains(t, err.Error(), "a@example.com")
	require.NotContains(t, err.Error(), "{{value}}")
}

func TestDuplicateWarningNote(t *testing.T) {
	t.Run("nothing matched, nothing to say", func(t *testing.T) {
		require.Empty(t, duplicateWarningNote(nil))
		require.Empty(t, duplicateWarningNote(&cmpTypes.RecordValueErrorSet{}))
	})

	t.Run("a soft match is named, with its value", func(t *testing.T) {
		dd := &cmpTypes.RecordValueErrorSet{Set: []cmpTypes.RecordValueError{{
			Kind:    "duplication_warning",
			Message: `The value "{{value}}" already exists in another record`,
			Meta:    map[string]interface{}{"field": "phone", "value": "+386 1 234"},
		}}}

		note := duplicateWarningNote(dd)
		require.Contains(t, note, "phone")
		require.Contains(t, note, "+386 1 234")
		require.Contains(t, note, "not strict")
	})

	t.Run("the same sentence for several records is said once", func(t *testing.T) {
		one := cmpTypes.RecordValueError{
			Kind:    "duplication_warning",
			Message: `The value "{{value}}" already exists in another record`,
			Meta:    map[string]interface{}{"field": "phone", "value": "+386 1 234"},
		}
		dd := &cmpTypes.RecordValueErrorSet{Set: []cmpTypes.RecordValueError{one, one, one}}

		require.Equal(t, 1, strings.Count(duplicateWarningNote(dd), "+386 1 234"))
	})
}
