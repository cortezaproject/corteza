package dal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentFormatterInit_noParams(t *testing.T) {
	IdentFormatter()
}

func TestIdentFormatterInit_okParams(t *testing.T) {
	IdentFormatter("k", "v")
}

func TestIdentFormatterInit_nokParams(t *testing.T) {
	assert.Panics(t, func() {
		IdentFormatter("k")
	})
}

func TestWithValidation(t *testing.T) {
	_, err := IdentFormatter().WithValidationE("true", nil)
	assert.NoError(t, err)
}

func TestWithValidation_nokExpr(t *testing.T) {
	_, err := IdentFormatter().WithValidationE("1variable", nil)
	assert.Error(t, err)
}

func TestWithValidation_nokExpr_panic(t *testing.T) {
	assert.Panics(t, func() {
		IdentFormatter().WithValidation("1variable", nil)
	})
}

func TestFormatting(t *testing.T) {
	tcc := []struct {
		name      string
		tpl       string
		params    []string
		validator string
		ident     string
		ok        bool
	}{
		{
			name:  "template without params",
			tpl:   "identifier",
			ident: "identifier",
			ok:    true,
		},
		{
			name:   "template with params",
			tpl:    "identifier_{{k}}",
			params: []string{"k", "v"},
			ident:  "identifier_v",
			ok:     true,
		},

		{
			name:      "template without params; validated ok",
			tpl:       "identifier",
			ident:     "identifier",
			validator: "true",
			ok:        true,
		},
		{
			name:      "template with params; validated ok",
			tpl:       "identifier_{{k}}",
			params:    []string{"k", "v"},
			ident:     "identifier_v",
			validator: "true",
			ok:        true,
		},

		{
			name:      "template without params; validated nok",
			tpl:       "identifier",
			ident:     "identifier",
			validator: "false",
			ok:        false,
		},
		{
			name:      "template with params; validated nok",
			tpl:       "identifier_{{k}}",
			params:    []string{"k", "v"},
			ident:     "identifier_v",
			validator: "false",
			ok:        false,
		},
	}

	ctx := context.Background()
	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			f := IdentFormatter()
			if c.validator != "" {
				f = f.WithValidation(c.validator, nil)
			}
			ident, ok := f.Format(ctx, c.tpl, c.params...)

			assert.Equal(t, c.ident, ident)
			assert.Equal(t, c.ok, ok)
		})
	}
}
