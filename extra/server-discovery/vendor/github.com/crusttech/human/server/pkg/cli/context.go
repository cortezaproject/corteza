package cli

import (
	"context"

	ctxwrap "github.com/SentimensRG/ctx"
	"github.com/SentimensRG/ctx/sigctx"
)

// Context is  small wrapper that returns sig-term bound context
//
// This can be used as (proper) background context that properly terminates
// all subroutines.
func Context() context.Context {
	return ctxwrap.AsContext(sigctx.New())
}
