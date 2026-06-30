package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

var dmlHandleRe = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

type (
	dmlSchemaReader interface {
		FindModels(ctx context.Context, connectionID uint64, f types.DmlModelFilter) ([]*types.DmlModel, error)
	}

	DmlMappingSvc struct {
		store  store.Storer
		schema dmlSchemaReader
		ac     dmlAC
	}
)

func NewDmlMappingSvc(s store.Storer, schema dmlSchemaReader, ac dmlAC) *DmlMappingSvc {
	return &DmlMappingSvc{store: s, schema: schema, ac: ac}
}

func (svc *DmlMappingSvc) Create(ctx context.Context, mp *types.DmlMapping) (*types.DmlMapping, error) {
	if !svc.ac.CanCreateDalConnection(ctx) {
		return nil, errors.Unauthorized("not allowed to manage DML mappings")
	}
	if mp.ID == 0 {
		mp.ID = id.Next()
	}
	if err := validateDmlMapping(mp); err != nil {
		return nil, err
	}
	mp.CreatedAt = time.Now()
	if err := store.CreateDmlMapping(ctx, svc.store, mp); err != nil {
		return nil, err
	}
	return mp, nil
}

func (svc *DmlMappingSvc) Find(ctx context.Context, f types.DmlMappingFilter) ([]*types.DmlMapping, error) {
	if !svc.ac.CanSearchDalConnections(ctx) {
		return nil, errors.Unauthorized("not allowed to access DML mappings")
	}
	set, _, err := store.SearchDmlMappings(ctx, svc.store, f)
	return set, err
}

func (svc *DmlMappingSvc) FindByID(ctx context.Context, mappingID uint64) (*types.DmlMapping, error) {
	if !svc.ac.CanSearchDalConnections(ctx) {
		return nil, errors.Unauthorized("not allowed to access DML mappings")
	}
	return svc.loadMapping(ctx, mappingID)
}

func (svc *DmlMappingSvc) FindByConnection(ctx context.Context, connectionID uint64) ([]*types.DmlMapping, error) {
	set, _, err := store.SearchDmlMappings(ctx, svc.store, types.DmlMappingFilter{
		ConnectionID: connectionID,
		Deleted:      filter.StateExcluded,
	})
	return set, err
}

func (svc *DmlMappingSvc) Update(ctx context.Context, mp *types.DmlMapping) (*types.DmlMapping, error) {
	if !svc.ac.CanCreateDalConnection(ctx) {
		return nil, errors.Unauthorized("not allowed to manage DML mappings")
	}
	if _, err := svc.loadMapping(ctx, mp.ID); err != nil {
		return nil, err
	}
	if err := validateDmlMapping(mp); err != nil {
		return nil, err
	}
	now := time.Now()
	mp.UpdatedAt = &now
	if err := store.UpdateDmlMapping(ctx, svc.store, mp); err != nil {
		return nil, err
	}
	return mp, nil
}

func (svc *DmlMappingSvc) DeleteByID(ctx context.Context, mappingID uint64) error {
	if !svc.ac.CanCreateDalConnection(ctx) {
		return errors.Unauthorized("not allowed to manage DML mappings")
	}
	mp, err := svc.loadMapping(ctx, mappingID)
	if err != nil {
		return err
	}
	now := time.Now()
	mp.DeletedAt = &now
	return store.UpdateDmlMapping(ctx, svc.store, mp)
}

// GenerateMapping introspects a connection and returns one suggested DmlMapping per external table.
// Callers may persist the returned mappings via Create.
func (svc *DmlMappingSvc) GenerateMapping(ctx context.Context, connectionID uint64) ([]*types.DmlMapping, error) {
	models, err := svc.schema.FindModels(ctx, connectionID, types.DmlModelFilter{})
	if err != nil {
		return nil, fmt.Errorf("dml: cannot introspect connection %d: %w", connectionID, err)
	}

	now := time.Now()
	out := make([]*types.DmlMapping, 0, len(models))
	for _, model := range models {
		mp := &types.DmlMapping{
			ID:           id.Next(),
			ConnectionID: connectionID,
			SourceIdent:  model.Ident,
			ModuleHandle: sanitizeDmlHandle(model.Ident),
			ModuleName:   model.Label,
			CreatedAt:    now,
		}
		for _, attr := range model.Attributes {
			mp.Columns = append(mp.Columns, &types.DmlColumnMap{
				SourceIdent: attr.Ident,
				FieldName:   attr.Ident,
				FieldKind:   dmlAttrTypeToFieldKind(attr.Type),
			})
		}
		out = append(out, mp)
	}
	return out, nil
}

func (svc *DmlMappingSvc) loadMapping(ctx context.Context, mappingID uint64) (*types.DmlMapping, error) {
	mp, err := store.LookupDmlMappingByID(ctx, svc.store, mappingID)
	if err != nil {
		return nil, err
	}
	if mp == nil || mp.DeletedAt != nil {
		return nil, fmt.Errorf("dml: mapping %d not found", mappingID)
	}
	return mp, nil
}

func validateDmlMapping(mp *types.DmlMapping) error {
	if mp.NamespaceHandle != "" && !dmlHandleRe.MatchString(mp.NamespaceHandle) {
		return fmt.Errorf("dml: invalid namespace handle %q", mp.NamespaceHandle)
	}
	if mp.Skip {
		return nil
	}
	if !dmlHandleRe.MatchString(mp.ModuleHandle) {
		return fmt.Errorf("dml: invalid module handle %q", mp.ModuleHandle)
	}
	hasCol := false
	for _, c := range mp.Columns {
		if !c.Skip {
			hasCol = true
			break
		}
	}
	if !hasCol {
		return fmt.Errorf("dml: mapping for %q has no non-skipped columns", mp.SourceIdent)
	}
	if mp.Identifier != "" {
		found := false
		for _, c := range mp.Columns {
			if c.SourceIdent == mp.Identifier && !c.Skip {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("dml: identifier column %q not found in mapping for %q", mp.Identifier, mp.SourceIdent)
		}
	}
	return nil
}

func sanitizeDmlHandle(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9_-]`).ReplaceAllString(s, "_")
	s = regexp.MustCompile(`_+`).ReplaceAllString(s, "_")
	s = strings.Trim(s, "_-")
	if len(s) == 0 || (s[0] >= '0' && s[0] <= '9') {
		s = "t_" + s
	}
	return s
}

func dmlAttrTypeToFieldKind(typeStr string) string {
	switch typeStr {
	case "corteza::dal:attribute-type:number":
		return "Number"
	case "corteza::dal:attribute-type:boolean":
		return "Bool"
	case "corteza::dal:attribute-type:timestamp",
		"corteza::dal:attribute-type:date",
		"corteza::dal:attribute-type:time":
		return "DateTime"
	case "corteza::dal:attribute-type:ref":
		return "Number"
	default:
		return "String"
	}
}
