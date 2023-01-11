package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
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
	return ctrl.svc.Evaluate(ctx, r.Expressions, r.Variables)
}
