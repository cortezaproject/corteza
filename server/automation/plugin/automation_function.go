package plugin

import (
	"context"

	"github.com/cortezaproject/corteza/server/automation/plugin/grpc"
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
)

type (
	cp struct {
		af grpc.AutomationFunction
	}
)

func (c *cp) Generate() *types.Function {
	var (
		f = c.af.Meta()
	)

	f.Handler = func(ctx context.Context, in *expr.Vars) (o *expr.Vars, err error) {
		out, err := c.af.Exec(ctx, in)

		if err != nil {
			return nil, err
		}

		o = &expr.Vars{}
		err = o.Assign(out)

		return
	}

	return f
}

func MakeAutomationFunction() *cp {
	return &cp{}
}

func (c *cp) SetPlugin(af grpc.AutomationFunction) {
	c.af = af
}
