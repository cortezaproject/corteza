package dml

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

var handleRe = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

// Mapping manages DmlMapping objects: auto-generation from introspection + CRUD.
type Mapping struct {
	conn  *Connection
	store store.Storer
}

func NewMapping(conn *Connection) *Mapping {
	return &Mapping{conn: conn, store: conn.store}
}

// GenerateMapping auto-generates a DmlMapping proposal from a live connection's schema.
func (m *Mapping) GenerateMapping(ctx context.Context, connectionID uint64) (*types.DmlMapping, error) {
	models, err := m.conn.FindModels(ctx, connectionID)
	if err != nil {
		return nil, fmt.Errorf("dml: cannot introspect connection %d: %w", connectionID, err)
	}

	mp := &types.DmlMapping{
		ID:           id.Next(),
		ConnectionID: connectionID,
	}

	for _, model := range models {
		tbl := &types.DmlTableMap{
			SourceIdent:  model.Ident,
			ModuleHandle: sanitizeHandle(model.Ident),
			ModuleName:   model.Label,
		}

		for _, attr := range model.Attributes {
			col := &types.DmlColumnMap{
				SourceIdent: attr.Ident,
				FieldName:   attr.Ident,
				FieldKind:   attrTypeToFieldKind(attr.Type),
			}
			tbl.Columns = append(tbl.Columns, col)
		}

		mp.Tables = append(mp.Tables, tbl)
	}

	return mp, nil
}

// Create persists a new mapping.
func (m *Mapping) Create(ctx context.Context, mp *types.DmlMapping) (*types.DmlMapping, error) {
	if mp.ID == 0 {
		mp.ID = id.Next()
	}

	if err := validateMapping(mp); err != nil {
		return nil, err
	}

	if err := saveMapping(ctx, m.store, mp); err != nil {
		return nil, err
	}

	return mp, nil
}

// FindByID returns a mapping by ID.
func (m *Mapping) FindByID(ctx context.Context, mappingID uint64) (*types.DmlMapping, error) {
	return loadMapping(ctx, m.store, mappingID)
}

// Update replaces a mapping.
func (m *Mapping) Update(ctx context.Context, mp *types.DmlMapping) (*types.DmlMapping, error) {
	if _, err := loadMapping(ctx, m.store, mp.ID); err != nil {
		return nil, err
	}

	if err := validateMapping(mp); err != nil {
		return nil, err
	}

	if err := saveMapping(ctx, m.store, mp); err != nil {
		return nil, err
	}

	return mp, nil
}

// DeleteByID removes a mapping.
func (m *Mapping) DeleteByID(ctx context.Context, mappingID uint64) error {
	return deleteMapping(ctx, m.store, mappingID)
}

// ----------------------------------------------------------------------------
// helpers
// ----------------------------------------------------------------------------

func sanitizeHandle(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9_-]`).ReplaceAllString(s, "_")
	s = regexp.MustCompile(`_+`).ReplaceAllString(s, "_")
	s = strings.Trim(s, "_-")
	if len(s) == 0 || (s[0] >= '0' && s[0] <= '9') {
		s = "t_" + s
	}
	return s
}

// attrTypeToFieldKind converts a DmlAttribute.Type string to a compose field kind.
func attrTypeToFieldKind(typeStr string) string {
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
		// Source foreign keys are opaque external IDs. A compose Record field
		// would reject them — ref values must point at an existing compose
		// record. Import the raw FK as a Number; remapping to compose IDs is a
		// future enhancement. (@todo optional ID-remap pass.)
		return "Number"
	default:
		return "String"
	}
}

func validateMapping(mp *types.DmlMapping) error {
	if mp.NamespaceHandle != "" && !handleRe.MatchString(mp.NamespaceHandle) {
		return fmt.Errorf("dml: invalid namespace handle %q", mp.NamespaceHandle)
	}

	handles := map[string]struct{}{}
	for _, t := range mp.Tables {
		if t.Skip {
			continue
		}
		if !handleRe.MatchString(t.ModuleHandle) {
			return fmt.Errorf("dml: invalid module handle %q for table %q", t.ModuleHandle, t.SourceIdent)
		}
		if _, dup := handles[t.ModuleHandle]; dup {
			return fmt.Errorf("dml: duplicate module handle %q", t.ModuleHandle)
		}
		handles[t.ModuleHandle] = struct{}{}

		hasCol := false
		for _, c := range t.Columns {
			if !c.Skip {
				hasCol = true
				break
			}
		}
		if !hasCol {
			return fmt.Errorf("dml: table %q has no non-skipped columns", t.SourceIdent)
		}

		if t.Identifier != "" {
			found := false
			for _, c := range t.Columns {
				if c.SourceIdent == t.Identifier && !c.Skip {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("dml: identifier column %q not found in table %q", t.Identifier, t.SourceIdent)
			}
		}
	}
	return nil
}
