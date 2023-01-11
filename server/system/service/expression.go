package service

import (
	"context"
	"fmt"

	"errors"

	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza/server/pkg/expr"
)

type (
	expression struct {
		parser gval.Language
	}
)

func Expression() *expression {
	return &expression{
		parser: expr.Parser(),
	}
}

func (svc *expression) Evaluate(ctx context.Context, expressions map[string]string, vars map[string]any) (out map[string]any, err error) {
	out = make(map[string]any, len(expressions))
	evalErrs := make(map[string]error, len(expressions))
	hasErr := false

	for k, e := range expressions {
		out[k], err = svc.parser.Evaluate(e, vars)
		evalErrs[k] = err
		hasErr = hasErr || err != nil
	}

	if !hasErr {
		return
	}

	err = nil
	errMsg := ""
	for k, e := range evalErrs {
		if e != nil {
			errMsg += fmt.Sprintf("%s: %s\n", k, e.Error())
		}
	}

	return nil, errors.New(errMsg)
}
