package app

import (
	"context"

	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/rbac"
)

func (app *HumanApp) InitExpr(ctx context.Context) (err error) {
	expr.Init(rbac.AllFunctions, expr.AllFunctions)

	return
}
