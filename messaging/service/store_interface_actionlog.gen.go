package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/actionlog.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
)

type (
	actionlogsStore interface {
		SearchActionlogs(ctx context.Context, f actionlog.Filter) (actionlog.ActionSet, actionlog.Filter, error)
		CreateActionlog(ctx context.Context, rr ...*actionlog.Action) error
		UpdateActionlog(ctx context.Context, rr ...*actionlog.Action) error
		PartialUpdateActionlog(ctx context.Context, onlyColumns []string, rr ...*actionlog.Action) error
		RemoveActionlog(ctx context.Context, rr ...*actionlog.Action) error
		RemoveActionlogByID(ctx context.Context, ID uint64) error

		TruncateActionlogs(ctx context.Context) error
	}
)
