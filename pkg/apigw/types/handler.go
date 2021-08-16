package types

import (
	"context"

	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	Execer interface {
		Exec(context.Context, *Scp) error
		Type() FilterKind
	}

	Sorter interface {
		Weight() int
	}

	ErrorHandler interface {
		Exec(context.Context, *Scp, error)
	}

	Stringer interface {
		String() string
	}

	WfExecer interface {
		Exec(ctx context.Context, workflowID uint64, p atypes.WorkflowExecParams) (*expr.Vars, atypes.Stacktrace, error)
	}

	Handler interface {
		Execer
		Stringer
		Sorter

		Merge([]byte) (Handler, error)
		Meta() FilterMeta
	}
)
