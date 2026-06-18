package dml

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx/types"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
	systemTypes "github.com/crusttech/human/server/system/types"
)

// DML persists its mappings and import runs as JSON blobs in the system
// settings KV table (owner 0). Access is whole-object by ID, so a dedicated
// table buys nothing; this avoids a codegen round-trip while still surviving
// restarts and being shared across nodes (unlike the previous in-memory maps).
//
// @todo promote to first-class store resources (dml_mapping / dml_import_run
//       .cue + `make codegen`) once these need filterable listing or RBAC.

const (
	settingMappingPrefix = "dml.mapping."
	settingRunPrefix     = "dml.run."
)

func mappingSettingName(id uint64) string { return fmt.Sprintf("%s%d", settingMappingPrefix, id) }
func runSettingName(id uint64) string     { return fmt.Sprintf("%s%d", settingRunPrefix, id) }

func saveMapping(ctx context.Context, s store.Storer, mp *systemTypes.DmlMapping) error {
	return putJSON(ctx, s, mappingSettingName(mp.ID), mp)
}

func loadMapping(ctx context.Context, s store.Storer, id uint64) (*systemTypes.DmlMapping, error) {
	mp := &systemTypes.DmlMapping{}
	ok, err := getJSON(ctx, s, mappingSettingName(id), mp)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("dml: mapping %d not found", id)
	}
	return mp, nil
}

func deleteMapping(ctx context.Context, s store.Storer, id uint64) error {
	return store.DeleteSettingValueByOwnedByName(ctx, s, 0, mappingSettingName(id))
}

func saveRun(ctx context.Context, s store.Storer, run *systemTypes.DmlImportRun) error {
	return putJSON(ctx, s, runSettingName(run.ID), run)
}

func loadRun(ctx context.Context, s store.Storer, id uint64) (*systemTypes.DmlImportRun, error) {
	run := &systemTypes.DmlImportRun{}
	ok, err := getJSON(ctx, s, runSettingName(id), run)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("dml: import run %d not found", id)
	}
	return run, nil
}

func putJSON(ctx context.Context, s store.Storer, name string, v any) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return store.UpsertSettingValue(ctx, s, &systemTypes.SettingValue{
		Name:    name,
		OwnedBy: 0,
		Value:   types.JSONText(raw),
	})
}

func getJSON(ctx context.Context, s store.Storer, name string, v any) (bool, error) {
	sv, err := store.LookupSettingValueByNameOwnedBy(ctx, s, name, 0)
	if errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if sv == nil {
		return false, nil
	}
	if err := json.Unmarshal([]byte(sv.Value), v); err != nil {
		return false, err
	}
	return true, nil
}
