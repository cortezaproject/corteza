package dal

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	identFormatter struct {
		identValidator       string
		identValidatorP      gval.Evaluable
		identValidatorParams map[string]any

		params []string
	}
)

// IdentFormatter returns an initialized ident formatter preconfigured with given
// base params.
//
// The ident formatter is primarily used for defining model and attribute identifiers.
//
// Base params are provided in key,value pairs and the function panics if an odd number of
// parameters is provided.
func IdentFormatter(baseParams ...string) identFormatter {
	out := identFormatter{}

	// Validate params; following what string replacer does
	err := out.validateFormatParams(baseParams...)
	if err != nil {
		panic(fmt.Sprintf("cannot initialize identFormatter: %s", err.Error()))
	}

	// Some preprocessing
	out.params = out.prepareFormatParams(baseParams...)
	return out
}

// WithValidation binds the given validation expression to the ident formatter
//
// The initial formatter remains unchanged.
func (f identFormatter) WithValidationE(validator string, params map[string]any) (_ identFormatter, err error) {
	f.identValidator = validator
	f.identValidatorP, err = expr.Parser().NewEvaluable(validator)
	f.identValidatorParams = params
	return f, err
}

// WithValidation binds the given validation expression to the ident formatter
//
// The initial formatter remains unchanged.
//
// The function panics if the expression can not be parsed.
func (f identFormatter) WithValidation(validator string, params map[string]any) identFormatter {
	out, err := f.WithValidationE(validator, params)
	if err != nil {
		panic(err)
	}
	return out
}

// Format returns a formatted identifier and a flag wether the resulting identifier is valid or not
//
// Parameters are provided in key,value pairs and the function panics if an odd number of
// parameters is provided.
func (f identFormatter) Format(ctx context.Context, template string, params ...string) (out string, ok bool) {
	var err error
	ok = true

	err = f.validateFormatParams(params...)
	if err != nil {
		panic(fmt.Sprintf("cannot format template: %s", err.Error()))
	}

	f.params = append(f.params, f.prepareFormatParams(params...)...)

	rpl := strings.NewReplacer(f.params...)
	out = rpl.Replace(template)

	if f.identValidator != "" {
		ok, err = f.identValidatorP.EvalBool(ctx, f.getEvalParams(out))
		ok = ok && (err == nil)
	}

	return
}

// getEvalParams is a helper to get a KV map of parameters for gval ident validation
func (f identFormatter) getEvalParams(ident string) (out map[string]any) {
	out = map[string]any{
		"ident": ident,
	}

	for k, v := range f.identValidatorParams {
		out[k] = v
	}

	return
}

func (f identFormatter) validateFormatParams(params ...string) error {
	if len(params)%2 == 1 {
		return errors.New("expecting even number of parameters")
	}

	return nil
}

func (f identFormatter) prepareFormatParams(params ...string) []string {
	for i := 0; i < len(params); i += 2 {
		params[i] = fmt.Sprintf("{{%s}}", params[i])
	}

	return params
}
