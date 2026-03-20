package types

import (
	"context"
	"errors"
	"fmt"

	pkgAst "github.com/cortezaproject/corteza/server/pkg/ast"
	execTypes "github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
)

type errorStep struct {
	message     *pkgAst.ASTNode
	recoverable bool
}

// ErrorStep constructs a step handler that throws an error when executed.
// If recoverable is true, the error is wrapped in RecoverableError, allowing
// the runtime to pause execution for user intervention.
// If recoverable is false, a plain error is returned and the catch walk begins.
func ErrorStep(message *pkgAst.ASTNode, recoverable bool) execTypes.StepHandler {
	return &errorStep{message: message, recoverable: recoverable}
}

func (s *errorStep) ExecN(ctx context.Context, r *execTypes.ExecRequest) (execTypes.ExecResponse, error) {
	msg := "error step triggered"
	if s.message != nil {
		// Evaluate the message expression against the scope.
		val, err := pkgAst.Eval(s.message, r.Scope)
		if err == nil && val != nil {
			if str := fmt.Sprintf("%v", val.Get()); str != "" {
				msg = str
			}
		}
	}

	if s.recoverable {
		return nil, execTypes.NewRecoverableError(msg)
	}
	return nil, errors.New(msg)
}
