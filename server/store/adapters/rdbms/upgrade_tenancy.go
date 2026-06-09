package rdbms

import (
	"context"
	"fmt"

	automationModel "github.com/crusttech/human/server/automation/model"
	federationModel "github.com/crusttech/human/server/federation/model"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store/adapters/rdbms/ddl"
	systemModel "github.com/crusttech/human/server/system/model"
)

// Multi-tenancy resource scoping migration.
//
// Adds the denormalised rel_tenant (and, for project-scoped resources,
// rel_project) columns plus a composite scope index to every scoped table.
//
// New deployments get these columns straight from the model definitions via
// createTablesFromModels; this fix only backfills the columns onto tables that
// already exist from an earlier schema. Columns default to 0, which the store
// scope guard treats as "unset" — existing data is reachable until a tenant is
// provisioned and the rows are backfilled (see multi-tenancy-resource-scoping.md).

type scopedTable struct {
	model        *dal.Model
	projectScope bool
}

// scopedTables enumerates every table that gained tenancy scope columns.
// The model carries the canonical rel_tenant / rel_project attribute
// definitions so the migration stays in lock-step with codegen.
func scopedTables() []scopedTable {
	return []scopedTable{
		// system — hierarchy tables (Tenant itself is scoped by `id`, not a
		// rel_tenant column, so it carries no scope guard here).
		{systemModel.Project, false},
		{systemModel.ProjectMember, true},
		{systemModel.TenantMembership, false},

		// system — tenant-scoped
		{systemModel.Application, false},
		{systemModel.AuthClient, false},
		{systemModel.LlmProvider, false},
		{systemModel.DalSensitivityLevel, false},
		{systemModel.Queue, false},
		{systemModel.ConfiguredConnection, false},
		{systemModel.UserGroup, false},
		{systemModel.ResourceTranslation, true},

		// system — project-scoped
		{systemModel.Agent, true},
		{systemModel.AiConversation, true},
		{systemModel.ApigwRoute, true},
		{systemModel.ApigwFilter, true},
		{systemModel.Chatbot, true},
		{systemModel.ChatbotSession, true},
		{systemModel.KnowledgeBase, true},
		{systemModel.Report, true},
		{systemModel.Template, true},
		{systemModel.Reminder, true},
		{systemModel.Attachment, true},
		{systemModel.Notification, true},
		{systemModel.DataPrivacyRequest, true},

		// automation — project-scoped
		{automationModel.Workflow, true},
		{automationModel.Trigger, true},
		{automationModel.Session, true},

		// federation
		{federationModel.Node, false},
		{federationModel.ExposedModule, true},
		{federationModel.SharedModule, true},
		{federationModel.ModuleMapping, true},
	}
}

// fix_2026_06_00_addTenancyScopeColumns backfills rel_tenant / rel_project
// columns and a composite scope index onto every pre-existing scoped table.
func fix_2026_06_00_addTenancyScopeColumns(ctx context.Context, s *Store) (err error) {
	for _, st := range scopedTables() {
		if err = addScopeColumns(ctx, s, st); err != nil {
			return fmt.Errorf("tenancy scope migration failed for %q: %w", st.model.Ident, err)
		}
	}
	return nil
}

func addScopeColumns(ctx context.Context, s *Store, st scopedTable) (err error) {
	table := st.model.Ident

	// only touch tables that already exist; brand-new tables get the columns
	// from the model at create time.
	if _, err = s.DataDefiner.TableLookup(ctx, table); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	tenantAttr := st.model.Attributes.FindByIdent("TenantID")
	if tenantAttr == nil {
		return fmt.Errorf("model %q is missing the TenantID attribute", table)
	}
	if err = addColumn(ctx, s, table, tenantAttr); err != nil {
		return err
	}

	indexFields := []*ddl.IndexField{{Column: "rel_tenant"}}

	if st.projectScope {
		projectAttr := st.model.Attributes.FindByIdent("ProjectID")
		if projectAttr == nil {
			return fmt.Errorf("model %q is missing the ProjectID attribute", table)
		}
		if err = addColumn(ctx, s, table, projectAttr); err != nil {
			return err
		}
		indexFields = append(indexFields, &ddl.IndexField{Column: "rel_project"})
	}

	return ensureScopeIndex(ctx, s, table, indexFields)
}

// ensureScopeIndex creates the composite tenancy index unless it already exists.
func ensureScopeIndex(ctx context.Context, s *Store, table string, fields []*ddl.IndexField) error {
	indexName := table + "_idxScope"

	idx, err := s.DataDefiner.IndexLookup(ctx, indexName, table)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if idx != nil {
		return nil
	}

	return s.DataDefiner.IndexCreate(ctx, table, &ddl.Index{
		TableIdent: table,
		Ident:      indexName,
		Type:       "BTREE",
		Fields:     fields,
	})
}
