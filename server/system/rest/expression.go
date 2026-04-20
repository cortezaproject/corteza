package rest

import (
	"context"
	"strings"

	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
)

type (
	Expression struct {
		svc exprService
	}

	exprService interface {
		Evaluate(context.Context, map[string]string, map[string]any) (map[string]any, error)
	}
)

func (Expression) New() *Expression {
	return &Expression{
		svc: service.DefaultExpression,
	}
}

func (ctrl *Expression) Evaluate(ctx context.Context, r *request.ExpressionEvaluate) (interface{}, error) {
	expressions := make(map[string]string, len(r.Expressions))
	for k, v := range r.Expressions {
		if v.Val != "" {
			expressions[k] = v.Val
		} else if len(v.Values) > 0 {
			expressions[k] = strings.Join(v.Values, ",")
		}
	}

	return ctrl.svc.Evaluate(ctx, expressions, r.Variables)
}
