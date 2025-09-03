package app

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
)

func (app *CortezaApp) initExpr(ctx context.Context) (err error) {
	expr.Init(rbac.AllFunctions)

	return
}
