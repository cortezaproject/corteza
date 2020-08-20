package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/actionlog.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
)

type (
	Actionlogs interface {
		SearchActionlogs(ctx context.Context, f actionlog.Filter) (actionlog.ActionSet, actionlog.Filter, error)

		CreateActionlog(ctx context.Context, rr ...*actionlog.Action) error

		TruncateActionlogs(ctx context.Context) error
	}
)

var _ *actionlog.Action
var _ context.Context

// SearchActionlogs returns all matching Actionlogs from store
func SearchActionlogs(ctx context.Context, s Actionlogs, f actionlog.Filter) (actionlog.ActionSet, actionlog.Filter, error) {
	return s.SearchActionlogs(ctx, f)
}

// CreateActionlog creates one or more Actionlogs in store
func CreateActionlog(ctx context.Context, s Actionlogs, rr ...*actionlog.Action) error {
	return s.CreateActionlog(ctx, rr...)
}

// TruncateActionlogs Deletes all Actionlogs from store
func TruncateActionlogs(ctx context.Context, s Actionlogs) error {
	return s.TruncateActionlogs(ctx)
}
