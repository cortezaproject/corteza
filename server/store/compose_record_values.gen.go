package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/compose_record_values.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeRecordValues interface {

		// Additional custom functions

		// ComposeRecordValueRefLookup (custom function)
		ComposeRecordValueRefLookup(ctx context.Context, _mod *types.Module, _field string, _ref uint64) (uint64, error)
	}
)

var _ *types.RecordValue
var _ context.Context

func ComposeRecordValueRefLookup(ctx context.Context, s ComposeRecordValues, _mod *types.Module, _field string, _ref uint64) (uint64, error) {
	return s.ComposeRecordValueRefLookup(ctx, _mod, _field, _ref)
}
