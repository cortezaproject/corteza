package dml

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	// dalReader is the slice of dal.FullService the connection service
	// consumes. Kept for the eventual real implementation; currently the
	// service is fixture-backed and never reaches the DAL.
	dalReader interface {
		GetConnectionByID(connectionID uint64) *dal.ConnectionWrap
		SearchModels(ctx context.Context) (dal.ModelSet, error)
		SearchConnectionIssues(connectionID uint64) []dal.Issue
	}

	// Connection projects pkg/dal connection + model state into API-facing
	// DML types. Today it returns canned fixtures; the real wiring against
	// store + dal lives behind the dalReader interface, ready to be
	// reconnected once the upstream surfaces are stable.
	Connection struct {
		store store.Storer
		dal   dalReader
	}
)

func NewConnection(s store.Storer, d dalReader) *Connection {
	return &Connection{store: s, dal: d}
}

// Find returns all DML connections matching the filter. Fixture data.
func (r *Connection) Find(_ context.Context, f types.DmlConnectionFilter) ([]*types.DmlConnection, error) {
	return filterConnections(fixtureConnections(), f), nil
}

// FindByID returns the DML connection with the given ID. Fixture data.
func (r *Connection) FindByID(_ context.Context, connectionID uint64) (*types.DmlConnection, error) {
	for _, c := range fixtureConnections() {
		if c.ID == connectionID {
			return c, nil
		}
	}
	return nil, fmt.Errorf("dml: connection %d not found", connectionID)
}

// FindModels returns the models attached to the given connection.
// connectionID 0 returns every model across every connection. Fixture data.
func (r *Connection) FindModels(_ context.Context, connectionID uint64) ([]*types.DmlModel, error) {
	out := fixtureModels()
	if connectionID == 0 {
		return out, nil
	}
	filtered := make([]*types.DmlModel, 0, len(out))
	for _, m := range out {
		if m.ConnectionID == connectionID {
			filtered = append(filtered, m)
		}
	}
	return filtered, nil
}

// FindModelByIdent resolves a single model on a connection by its ident.
// Fixture data.
func (r *Connection) FindModelByIdent(_ context.Context, connectionID uint64, ident string) (*types.DmlModel, error) {
	for _, m := range fixtureModels() {
		if m.ConnectionID == connectionID && m.Ident == ident {
			return m, nil
		}
	}
	return nil, fmt.Errorf("dml: model %q not found on connection %d", ident, connectionID)
}

// ----------------------------------------------------------------------------
// Fixtures
//
// Temporary stubs so the API surface can be exercised before the DAL
// projection is wired. Replace fixtureConnections / fixtureModels with
// calls into r.store + r.dal once the real plumbing lands.
// ----------------------------------------------------------------------------

func fixtureConnections() []*types.DmlConnection {
	return []*types.DmlConnection{
		{
			ID:           1,
			Handle:       "primary-database",
			Type:         "corteza::system:primary-dal-connection",
			Label:        "Primary database",
			Driver:       "corteza::dal:driver:rdbms",
			Capabilities: []string{"create", "update", "delete", "search", "lookup"},
			ModelIdent:   "{{namespace}}_{{module}}",
		},
		{
			ID:           2,
			Handle:       "warehouse",
			Type:         "corteza::system:dal-connection",
			Label:        "Analytics warehouse",
			Driver:       "corteza::dal:driver:rdbms",
			Capabilities: []string{"create", "search", "lookup"},
			ModelIdent:   "wh_{{module}}",
		},
		{
			ID:           3,
			Handle:       "crm-api",
			Type:         "corteza::system:dal-connection",
			Label:        "External CRM (REST)",
			Driver:       "corteza::dal:driver:api",
			Capabilities: []string{"search", "lookup"},
		},
	}
}

func fixtureModels() []*types.DmlModel {
	const (
		codecAlias = "corteza::dal:attribute-codec:alias"
		typeID     = "corteza::dal:attribute-type:id"
		typeRef    = "corteza::dal:attribute-type:ref"
		typeText   = "corteza::dal:attribute-type:text"
		typeJSON   = "corteza::dal:attribute-type:json"
		typeTS     = "corteza::dal:attribute-type:timestamp"
		typeBool   = "corteza::dal:attribute-type:boolean"
	)

	return []*types.DmlModel{
		{
			ConnectionID: 1,
			Ident:        "compose_record",
			Label:        "Compose record",
			ResourceType: "corteza::compose:record",
			Attributes: []*types.DmlAttribute{
				{Ident: "ID", PrimaryKey: true, Type: typeID, Store: codecAlias},
				{Ident: "ModuleID", Type: typeRef, Store: codecAlias},
				{Ident: "NamespaceID", Type: typeRef, Store: codecAlias},
				{Ident: "Values", Type: typeJSON, Store: codecAlias},
				{Ident: "Meta", Type: typeJSON, Store: codecAlias},
				{Ident: "CreatedAt", Sortable: true, Type: typeTS, Store: codecAlias},
				{Ident: "UpdatedAt", Sortable: true, Type: typeTS, Store: codecAlias},
				{Ident: "DeletedAt", Sortable: true, Type: typeTS, Store: codecAlias},
				{Ident: "OwnedBy", Type: typeRef, Store: codecAlias},
				{Ident: "CreatedBy", Type: typeRef, Store: codecAlias},
				{Ident: "UpdatedBy", Type: typeRef, Store: codecAlias},
				{Ident: "DeletedBy", Type: typeRef, Store: codecAlias},
			},
		},
		{
			ConnectionID: 1,
			Ident:        "users",
			Label:        "User",
			ResourceType: "corteza::system:user",
			Attributes: []*types.DmlAttribute{
				{Ident: "ID", PrimaryKey: true, Type: typeID, Store: codecAlias},
				{Ident: "Email", Sortable: true, Type: typeText, Store: codecAlias},
				{Ident: "EmailConfirmed", Type: typeBool, Store: codecAlias},
				{Ident: "Username", Sortable: true, Type: typeText, Store: codecAlias},
				{Ident: "Name", Sortable: true, Type: typeText, Store: codecAlias},
				{Ident: "Handle", Type: typeText, Store: codecAlias},
				{Ident: "Kind", Sortable: true, Type: typeText, Store: codecAlias},
				{Ident: "Meta", Type: typeJSON, Store: codecAlias},
				{Ident: "SuspendedAt", Sortable: true, Type: typeTS, Store: codecAlias},
				{Ident: "CreatedAt", Sortable: true, Type: typeTS, Store: codecAlias},
				{Ident: "UpdatedAt", Sortable: true, Type: typeTS, Store: codecAlias},
				{Ident: "DeletedAt", Sortable: true, Type: typeTS, Store: codecAlias},
			},
		},
		{
			ConnectionID: 2,
			Ident:        "wh_record_fact",
			Label:        "Record fact (warehouse)",
			ResourceType: "corteza::warehouse:record-fact",
			Attributes: []*types.DmlAttribute{
				{Ident: "ID", PrimaryKey: true, Type: typeID, Store: codecAlias},
				{Ident: "RecordID", Type: typeRef, Store: codecAlias},
				{Ident: "IngestedAt", Sortable: true, Type: typeTS, Store: codecAlias},
				{Ident: "Payload", Type: typeJSON, Store: codecAlias},
			},
		},
		{
			ConnectionID: 3,
			Ident:        "crm_contact",
			Label:        "CRM contact",
			ResourceType: "corteza::crm:contact",
			Attributes: []*types.DmlAttribute{
				{Ident: "ID", PrimaryKey: true, Type: typeID, Store: codecAlias},
				{Ident: "Email", Filterable: true, Type: typeText, Store: codecAlias},
				{Ident: "FirstName", Type: typeText, Store: codecAlias},
				{Ident: "LastName", Type: typeText, Store: codecAlias},
				{Ident: "LifecycleStage", Type: typeText, Store: codecAlias},
			},
		},
	}
}

func filterConnections(in []*types.DmlConnection, f types.DmlConnectionFilter) []*types.DmlConnection {
	if len(f.ConnectionID) == 0 && f.Handle == "" && f.Type == "" {
		return in
	}
	ids := map[string]struct{}{}
	for _, id := range f.ConnectionID {
		if id != "" {
			ids[id] = struct{}{}
		}
	}
	out := make([]*types.DmlConnection, 0, len(in))
	for _, c := range in {
		if f.Handle != "" && c.Handle != f.Handle {
			continue
		}
		if f.Type != "" && c.Type != f.Type {
			continue
		}
		if len(ids) > 0 {
			if _, ok := ids[fmt.Sprintf("%d", c.ID)]; !ok {
				continue
			}
		}
		out = append(out, c)
	}
	return out
}
