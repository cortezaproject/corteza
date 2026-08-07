package rdbms

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/crusttech/human/server/compose/model"
	"github.com/crusttech/human/server/compose/types"
	discovery "github.com/crusttech/human/server/discovery/types"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/filter"
	labelsType "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/ddl"
	systemModel "github.com/crusttech/human/server/system/model"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// RDBMS database fixes
//
// Schema changes that can not be automatically applied or complex changes
// that require some logic to be applied are handled here.
//
// Function names should start with "fix" and version.
// This does not have any effect on how fixes are executed, only for organisation purposes.

var (
	// all enabled fix function need to be listed here
	fixesPre = []func(context.Context, *Store) error{
		// Scope columns must exist before any fix below queries a scoped table:
		// the generated queries unconditionally select rel_tenant/rel_project,
		// so e.g. migrateOldComposeRecordValues would fail on a pre-tenancy DB.
		// No-ops on fresh databases (tables don't exist yet; they're created
		// from the models, columns included).
		fix_2026_06_00_addTenancyScopeColumns,
		fix_2026_05_00_addCreatedByAgentToComposeResources,

		fix_2022_09_00_extendComposeModuleForPrivacyAndDAL,
		fix_2022_09_00_extendComposeModuleFieldsForPrivacyAndDAL,
		fix_2022_09_00_dropObsoleteComposeModuleFields,
		fix_2022_09_00_extendDalConnectionsForMeta,
		fix_2022_09_00_renameModuleColOnComposeRecords,
		fix_2022_09_00_addValuesOnComposeRecords,
		fix_2022_09_00_migrateOldComposeRecordValues,
		fix_2022_09_00_addRevisionOnComposeRecords,
		fix_2022_09_00_addMetaOnComposeRecords,
		fix_2022_09_00_addMissingNodeIdOnFederationMapping,
		fix_2023_03_00_migrateComposeModuleConfigForRecordDeDup,
		fix_2022_09_07_changePostgresIdColumnsDatatype,
		fix_2022_09_00_migrateComposeModuleDiscoveryConfigSettings,
		fix_2023_03_00_migrateComposePageMeta,
		fix_2024_09_03_dropFederationNodeSyncPrimaryKey,
		fix_2024_09_03_renameFederationNodeSyncNodeID,
		fix_2024_09_03_renameFederationNodeSyncComposeID,
		fix_2024_09_03_addFederationNodeSyncNodeIDIndex,
		fix_2024_9_7_migrateLabelsValueToJsonb,
		fix_2026_06_29_addDeletedAtOnDmlTables,
		fix_2026_06_29_addSourceIdentOnDmlMappings,
		fix_2026_07_00_addRevisionColumnsOnProjects,
		fix_2026_07_14_addModeOnProjects,
		fix_2026_07_21_dropModeOnProjects,
	}

	fixesPost = append([]func(context.Context, *Store) error{
		fix_2024_09_05_addRelResourceRoleMembershipColumn,
		fix_2024_09_05_addUserGroupReferenceToUser,
		fix_2026_04_00_addChatbotColumnToAgents,
		fix_2026_05_00_addSourceOnConnections,
		fix_2026_05_00_addStateOnChatbotSessions,
		fix_2026_06_00_addTenancyScopeColumns,
		fix_2026_07_28_addRelRevisionOnProjectWorkItems,
		fix_2026_07_30_dropGlobalUniqueHandleOnAgents,
		fix_2026_07_30_dropGlobalUniqueHandleOnChatbots,
		fix_2026_07_30_addArchivedAtOnProjects,
		fix_2026_07_30_backfillProjectRefOnComposeResources,
		fix_2026_07_31_addApprovalOnProjects,
		fix_2026_08_06_addCreatedByOnComposeNamespace,
	}, actionlogFixes...)

	// actionlog-only additive column fixes. Shared here so both the main Upgrade
	// and the standalone UpgradeActionlog (dedicated actionlog DB) apply them.
	actionlogFixes = []func(context.Context, *Store) error{
		fix_2026_06_22_addDeltaOnActionlog,
		fix_2026_06_22_addOldStateOnActionlog,
		fix_2026_07_13_addRelProjectOnActionlog,
		fix_2026_07_14_addRelResourceProjectOnActionlog,
		fix_2026_07_14_addRelTenantOnActionlog,
		fix_2026_07_14_addRelRootProjectOnActionlog,
		fix_2026_07_14_makeActionlogJsonColsNullable,
		fix_2026_07_28_addRelResourceRevisionOnActionlog,
	}
)

func fix_2026_06_22_addDeltaOnActionlog(ctx context.Context, s *Store) error {
	return addColumn(ctx, s,
		"actionlog",
		&dal.Attribute{
			Ident: "Delta",
			Type:  &dal.TypeJSON{Nullable: true},
			Store: &dal.CodecAlias{Ident: "delta"},
		},
	)
}

func fix_2026_06_22_addOldStateOnActionlog(ctx context.Context, s *Store) error {
	return addColumn(ctx, s,
		"actionlog",
		&dal.Attribute{
			Ident: "OldState",
			Type:  &dal.TypeJSON{Nullable: true},
			Store: &dal.CodecAlias{Ident: "old_state"},
		},
	)
}

// Adds the rel_project column (+ index) that scopes action-log events to a
// project. Reuses the generated model attribute so the column type matches the
// canonical ProjectRefField. Additive & idempotent — addColumn/index no-op if
// already present, and skip cleanly when the table doesn't exist yet.
//
// No backfill: a project id is only known from the request scope at write time
// (see actionlog enrich()); it can't be recovered from existing rows, so old
// events stay at project 0.
func fix_2026_07_13_addRelProjectOnActionlog(ctx context.Context, s *Store) error {
	if _, err := s.DataDefiner.TableLookup(ctx, "actionlog"); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	attr := systemModel.Action.Attributes.FindByIdent("ProjectID")
	if attr == nil {
		return fmt.Errorf("actionlog model is missing the ProjectID attribute")
	}
	if err := addColumn(ctx, s, "actionlog", attr); err != nil {
		return err
	}

	const indexName = "actionlog_rel_project"
	idx, err := s.DataDefiner.IndexLookup(ctx, indexName, "actionlog")
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if idx != nil {
		return nil
	}

	return s.DataDefiner.IndexCreate(ctx, "actionlog", &ddl.Index{
		TableIdent: "actionlog",
		Ident:      indexName,
		Type:       "BTREE",
		Fields:     []*ddl.IndexField{{Column: "rel_project"}},
	})
}

func fix_2026_07_14_addRelResourceProjectOnActionlog(ctx context.Context, s *Store) error {
	if _, err := s.DataDefiner.TableLookup(ctx, "actionlog"); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	attr := systemModel.Action.Attributes.FindByIdent("ResourceProjectID")
	if attr == nil {
		return fmt.Errorf("actionlog model is missing the ResourceProjectID attribute")
	}
	if err := addColumn(ctx, s, "actionlog", attr); err != nil {
		return err
	}

	const indexName = "actionlog_rel_resource_project"
	idx, err := s.DataDefiner.IndexLookup(ctx, indexName, "actionlog")
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if idx != nil {
		return nil
	}

	return s.DataDefiner.IndexCreate(ctx, "actionlog", &ddl.Index{
		TableIdent: "actionlog",
		Ident:      indexName,
		Type:       "BTREE",
		Fields:     []*ddl.IndexField{{Column: "rel_resource_project"}},
	})
}

func fix_2026_07_14_addRelTenantOnActionlog(ctx context.Context, s *Store) error {
	if _, err := s.DataDefiner.TableLookup(ctx, "actionlog"); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	attr := systemModel.Action.Attributes.FindByIdent("TenantID")
	if attr == nil {
		return fmt.Errorf("actionlog model is missing the TenantID attribute")
	}
	return addColumn(ctx, s, "actionlog", attr)
}

func fix_2026_07_14_addRelRootProjectOnActionlog(ctx context.Context, s *Store) error {
	if _, err := s.DataDefiner.TableLookup(ctx, "actionlog"); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	attr := systemModel.Action.Attributes.FindByIdent("RootProjectID")
	if attr == nil {
		return fmt.Errorf("actionlog model is missing the RootProjectID attribute")
	}
	if err := addColumn(ctx, s, "actionlog", attr); err != nil {
		return err
	}

	const indexName = "actionlog_rel_root_project"
	idx, err := s.DataDefiner.IndexLookup(ctx, indexName, "actionlog")
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if idx != nil {
		return nil
	}

	return s.DataDefiner.IndexCreate(ctx, "actionlog", &ddl.Index{
		TableIdent: "actionlog",
		Ident:      indexName,
		Type:       "BTREE",
		Fields:     []*ddl.IndexField{{Column: "rel_root_project"}},
	})
}

// Adds the rel_resource_revision column (+ index) that attributes an
// action-log event to the revision owning (or assigned to) the affected
// resource. Reuses the generated model attribute so the column type matches
// the canonical attribute. Additive & idempotent — addColumn/index no-op if
// already present, and skip cleanly when the table doesn't exist yet.
//
// No backfill: the revision is only known from the resource at write time
// (see actionlog enrich() and the per-resource ToAction() implementations);
// it can't be recovered for existing rows, so old events stay unattributed
// (0), same as fix_2026_07_14_addRelResourceProjectOnActionlog.
func fix_2026_07_28_addRelResourceRevisionOnActionlog(ctx context.Context, s *Store) error {
	if _, err := s.DataDefiner.TableLookup(ctx, "actionlog"); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	attr := systemModel.Action.Attributes.FindByIdent("ResourceRevisionID")
	if attr == nil {
		return fmt.Errorf("actionlog model is missing the ResourceRevisionID attribute")
	}
	if err := addColumn(ctx, s, "actionlog", attr); err != nil {
		return err
	}

	const indexName = "actionlog_rel_resource_revision"
	idx, err := s.DataDefiner.IndexLookup(ctx, indexName, "actionlog")
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if idx != nil {
		return nil
	}

	return s.DataDefiner.IndexCreate(ctx, "actionlog", &ddl.Index{
		TableIdent: "actionlog",
		Ident:      indexName,
		Type:       "BTREE",
		Fields:     []*ddl.IndexField{{Column: "rel_resource_revision"}},
	})
}

// Drops NOT NULL from delta and old_state on actionlog. When the actionlog
// table is freshly created from the model (UpgradeActionlog on an empty DB),
// those columns are created as NOT NULL before the model gained Nullable:true.
func fix_2026_07_14_makeActionlogJsonColsNullable(ctx context.Context, s *Store) error {
	tbl, err := s.DataDefiner.TableLookup(ctx, "actionlog")
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	driver := s.DB.DriverName()
	for _, col := range []string{"delta", "old_state"} {
		c := tbl.ColumnByIdent(col)
		if c == nil || c.Type.Null {
			continue
		}

		var query string
		switch {
		case strings.HasPrefix(driver, "postgres"):
			query = fmt.Sprintf("ALTER TABLE actionlog ALTER COLUMN %s DROP NOT NULL", col)
		case strings.HasPrefix(driver, "mysql"):
			query = fmt.Sprintf("ALTER TABLE actionlog MODIFY COLUMN %s JSON NULL", col)
		default:
			continue
		}

		if _, err = s.DB.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("failed to make actionlog.%s nullable: %w", col, err)
		}
	}

	return nil
}

func fix_2026_06_29_addSourceIdentOnDmlMappings(ctx context.Context, s *Store) error {
	return addColumn(ctx, s, "dml_mappings", &dal.Attribute{
		Ident: "SourceIdent",
		Type:  &dal.TypeText{Length: 256},
		Store: &dal.CodecAlias{Ident: "source_ident"},
	})
}

func fix_2026_06_29_addDeletedAtOnDmlTables(ctx context.Context, s *Store) error {
	attr := &dal.Attribute{
		Ident: "DeletedAt",
		Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
		Store: &dal.CodecAlias{Ident: "deleted_at"},
	}
	for _, table := range []string{"dml_connections", "dml_mappings", "dml_import_runs"} {
		if err := addColumn(ctx, s, table, attr); err != nil {
			return err
		}
	}
	return nil
}

// fix_2026_05_00_addStateOnChatbotSessions backfills the JSON `state` column
// introduced alongside the unified chatbot session model. Default is `[]`
// because the underlying Go type (types.ChatbotSessionState) is a slice — a
// `{}` backfill would fail to scan on first read.
func fix_2026_05_00_addStateOnChatbotSessions(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"chatbot_sessions",
		&dal.Attribute{
			Ident: "State",
			Type:  &dal.TypeJSON{DefaultValue: "[]"},
			Store: &dal.CodecAlias{Ident: "state"},
		},
	)
}

func fix_2026_04_00_addChatbotColumnToAgents(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"agents",
		&dal.Attribute{
			Ident: "Chatbot",
			Type:  &dal.TypeJSON{DefaultValue: "{}"},
			Store: &dal.CodecAlias{Ident: "chatbot"},
		},
	)
}

func fix_2026_05_00_addSourceOnConnections(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"connections",
		&dal.Attribute{Ident: "source", Type: &dal.TypeText{Length: 1024, HasDefault: true, DefaultValue: ""}},
	)
}

// Namespaces record their author so cleanup and auditing can ask who made a
// resource without reading it out of the slug. Existing rows keep 0 — the
// information was never stored, so there is nothing to backfill from.
func fix_2026_08_06_addCreatedByOnComposeNamespace(ctx context.Context, s *Store) error {
	return addColumn(ctx, s, model.Namespace.Ident, model.Namespace.Attributes.FindByIdent("CreatedBy"))
}

func fix_2026_05_00_addCreatedByAgentToComposeResources(ctx context.Context, s *Store) (err error) {
	for _, m := range []*dal.Model{model.Namespace, model.Module, model.ModuleField, model.Record, model.Page, model.PageLayout} {
		if err = addColumn(ctx, s, m.Ident, m.Attributes.FindByIdent("CreatedByAgent")); err != nil {
			return err
		}
	}
	return nil
}

func fix_2022_09_00_migrateComposeModuleDiscoveryConfigSettings(ctx context.Context, s *Store) (err error) {
	type (
		oldS struct {
			Discovery struct {
				Public    interface{} `json:"public"`
				Private   interface{} `json:"private"`
				Protected interface{} `json:"protected"`
			} `json:"discovery"`

			DAL             interface{} `json:"dal"`
			Privacy         interface{} `json:"privacy"`
			RecordRevisions interface{} `json:"recordRevisions"`
			RecordDeDup     interface{} `json:"recordDeDup"`
		}

		result struct {
			Result discovery.Result `json:"result"`
		}

		updateModule struct {
			ID     uint64 `json:"id"`
			Config []byte `json:"config"`
		}
	)

	var (
		log   = s.log(ctx)
		query string
		auxID uint64
		aux   []byte
		rows  *sql.Rows
		ss    oldS

		uu []updateModule

		driver = s.DB.DriverName()
	)

	const (
		getModuleDiscoverySettings = `
			SELECT id, compose_module.config AS discovery
			FROM compose_module`

		// The document is bound as a parameter rather than inlined into a CAST.
		// sqlite has no JSON type, so CAST(<document> AS JSON) falls through to
		// NUMERIC affinity and stores the integer 0 -- and this fix runs on
		// every boot, so it was wiping the config of every module in the
		// database. The damage is not even self-contained: a module whose
		// config reads back as 0 can no longer be scanned, so the next fix in
		// the list (migrateComposeModuleConfigForRecordDeDup, which loads all
		// modules through the store) fails and takes the whole upgrade with it.
		// The column is JSON on mysql and TEXT on sqlite; both accept the
		// document as a plain string parameter.
		updateModuleDiscoverySettings = `
			UPDATE compose_module
			SET config = ?
			WHERE id = ?`

		// Postgres keeps an explicit cast: jsonb will not take a text
		// parameter without one.
		updatePSQLModuleDiscoverySettings = `
			UPDATE compose_module
			SET config = $1::jsonb
			WHERE id = $2`
	)

	// 1. Check if module has discovery settings
	// 2. If yes, migrate them to new format from json to json slice for multiple lang support
	// 3. Save module

	_, err = s.DataDefiner.TableLookup(ctx, model.Module.Ident)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Debug("skipping module config recordDeDup migration: compose_module table not found")
			return nil
		}
		return err
	}

	_ = s.Tx(ctx, func(ctx context.Context, s store.Storer) (err error) {
		query = fmt.Sprintf(getModuleDiscoverySettings)
		rows, err = s.(*Store).DB.QueryContext(ctx, query)
		if err != nil {
			return
		}

		defer func() {
			// assign error to return value...
			err = rows.Close()
		}()

		for rows.Next() {
			if err = rows.Err(); err != nil {
				return
			}

			err = rows.Scan(&auxID, &aux)
			if err != nil {
				continue
			}

			if aux == nil {
				continue
			}

			err = json.Unmarshal(aux, &ss)
			if err != nil {
				continue
			}

			var (
				bb             []byte
				settings       discovery.ModuleMeta
				migrateSetting = func(input interface{}) (out result) {
					out = result{
						Result: discovery.Result{
							Lang:   "",
							Fields: []string{},
						},
					}
					var (
						ok  bool
						ii  map[string]interface{}
						rr  []interface{}
						res interface{}
					)

					if input != nil {
						ii, ok = input.(map[string]interface{})
						if ok {
							if ii["result"] != nil {
								rr, ok = ii["result"].([]interface{})
								if !ok {
									res, ok = ii["result"].(interface{})
									if ok {
										rr = append(rr, res)
									}
								}
								if ok {
									for _, r := range rr {
										out.Result.Lang = r.(map[string]interface{})["lang"].(string)
										out.Result.Fields = []string{}
										if _, ok = r.(map[string]interface{})["fields"].([]interface{}); ok {
											for _, f := range r.(map[string]interface{})["fields"].([]interface{}) {
												out.Result.Fields = append(out.Result.Fields, f.(string))
											}
										}
									}
								}
							}
						}
					}
					return
				}
			)

			settings.Public.Result = append(settings.Public.Result, migrateSetting(ss.Discovery.Public).Result)
			settings.Private.Result = append(settings.Private.Result, migrateSetting(ss.Discovery.Private).Result)
			settings.Protected.Result = append(settings.Protected.Result, migrateSetting(ss.Discovery.Protected).Result)

			ss.Discovery.Public = settings.Public
			ss.Discovery.Private = settings.Private
			ss.Discovery.Protected = settings.Protected

			bb, err = json.Marshal(ss)
			if err != nil {
				continue
			}

			uu = append(uu, updateModule{
				ID:     auxID,
				Config: bb,
			})
		}

		for _, u := range uu {
			if driver == "postgres" || driver == "postgres+debug" {
				query = updatePSQLModuleDiscoverySettings
			} else {
				query = s.(*Store).DB.Rebind(updateModuleDiscoverySettings)
			}
			log.Debug("saving migrated module.config.discovery settings", logger.Uint64("id", u.ID))
			_, err = s.(*Store).DB.ExecContext(ctx, query, string(u.Config), u.ID)
			if err != nil {
				log.Debug("error saving migrated module.config.discovery settings", logger.Uint64("id", u.ID))
				continue
			}
		}

		return
	})

	return
}

func fix_2023_03_00_migrateComposePageMeta(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"compose_page",
		model.Page.Attributes.FindByIdent("meta"),
	)
}

func fix_2022_09_00_extendComposeModuleForPrivacyAndDAL(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"compose_module",
		&dal.Attribute{Ident: "config", Type: &dal.TypeJSON{DefaultValue: "{}"}},
	)
}

func fix_2022_09_00_extendComposeModuleFieldsForPrivacyAndDAL(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"compose_module_field",
		&dal.Attribute{Ident: "config", Type: &dal.TypeJSON{DefaultValue: "{}"}},
	)
}

func fix_2022_09_00_dropObsoleteComposeModuleFields(ctx context.Context, s *Store) (err error) {
	return dropColumns(ctx, s,
		"compose_module_field",
		"is_private",
		"is_visible",
	)
}

func fix_2022_09_00_extendDalConnectionsForMeta(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"dal_connections",
		&dal.Attribute{Ident: "meta", Type: &dal.TypeJSON{DefaultValue: "{}"}},
	)
}

func fix_2022_09_00_renameModuleColOnComposeRecords(ctx context.Context, s *Store) (err error) {
	return renameColumn(ctx, s, "compose_record", "module_id", "rel_module")
}

func fix_2022_09_00_addValuesOnComposeRecords(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"compose_record",
		model.Record.Attributes.FindByIdent("Values"),
	)
}

func fix_2022_09_00_migrateOldComposeRecordValues(ctx context.Context, s *Store) (err error) {
	var (
		// default value, should be enough for most cases
		// we'll allow users to override this value if they need to
		// with UPGRADE_MIGRATE_OLD_COMPOSE_RECORD_VALUES_BATCH_SIZE env var
		recordSliceSize = 1000

		log = s.log(ctx)
	)

	_, err = s.DataDefiner.TableLookup(ctx, model.Record.Ident)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err

	}

	// parse UPGRADE_MIGRATE_OLD_COMPOSE_RECORD_VALUES_BATCH_SIZE and set value to recordSliceSize if valid
	if aux, set := os.LookupEnv("UPGRADE_MIGRATE_OLD_COMPOSE_RECORD_VALUES_BATCH_SIZE"); set {
		if auxInt := cast.ToInt(aux); auxInt > 0 {
			recordSliceSize = auxInt
		}
	}

	const (
		crvTableIdent = "compose_record_value"

		recordsPerModule = `
	SELECT id
	  FROM compose_record
	 WHERE rel_namespace = %d AND rel_module = %d AND id > %d AND deleted_at IS NULL ORDER BY id LIMIT %d
`

		// used with sprintf and 2 queries because of some limitation in pg driver & percona
		//
		// using subquery does not work:
		// This version of MySQL doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
		recValuesPerModule = `
	SELECT record_id, name, value, ref, place
	  FROM compose_record_value
	 WHERE record_id IN (%s)
	   AND deleted_at IS NULL
	 ORDER BY record_id, name, place`
	)

	// check if old record-value table exists
	_, err = s.DataDefiner.TableLookup(ctx, crvTableIdent)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Debug("skipping record value migration: compose_record_values table not found, " +
				"all record values migrated from 2022.3 format (compose_record_values table) " +
				"to 2022.9 format (values column on compose_record table)")
			return nil
		}
		return err
	}

	var (
		query     string
		recordIDs []string

		modules types.ModuleSet
		fields  types.ModuleFieldSet
		field   *types.ModuleField
		rows    *sql.Rows

		sliceLastRecordID uint64

		recordID, ref uint64
		place         uint
		value, name   string

		values map[uint64]map[string][]any
		intVal any

		totalRecords = count(ctx, s, model.Record.Ident)
		countRecords = 0
	)

	// parse UPGRADE_MIGRATE_OLD_COMPOSE_RECORD_VALUES_BATCH_SIZE and set value
	// to recordSliceSize if valid
	if aux, set := os.LookupEnv("UPGRADE_MIGRATE_OLD_COMPOSE_RECORD_VALUES_BATCH_SIZE"); set {
		if auxInt := cast.ToInt(aux); auxInt > 0 {
			recordSliceSize = auxInt
		}
	}

	modules, _, err = s.SearchComposeModules(ctx, types.ModuleFilter{Deleted: filter.StateInclusive})
	if err != nil {
		return
	}

	log.Info(
		"preparing to migrate record values",
		zap.Int("modules", len(modules)),
		zap.Int("records", totalRecords),
		zap.Int("batch-size", recordSliceSize),
	)

	// iterate through modules
	for _, mod := range modules {
		fields, _, err = s.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
		if err != nil {
			return
		}

		perModLog := log.With(
			zap.String("handle", mod.Handle),
			logger.Uint64("id", mod.ID),
		)

		err = func() (err error) {
			sliceLastRecordID = 0

			for {
				bmStart := time.Now()
				values = make(map[uint64]map[string][]any, recordSliceSize)
				recordIDs = make([]string, 0, recordSliceSize)

				err = func() (err error) {
					query = fmt.Sprintf(recordsPerModule, mod.NamespaceID, mod.ID, sliceLastRecordID, recordSliceSize)
					// println(query)
					rows, err = s.DB.QueryContext(ctx, query)
					if err != nil {
						return
					}

					defer func() {
						// assign error to return value...
						err = rows.Close()
					}()

					for rows.Next() {
						if err = rows.Err(); err != nil {
							return
						}

						err = rows.Scan(&value)
						if err != nil {
							return
						}

						recordIDs = append(recordIDs, value)
					}

					if len(recordIDs) == 0 {
						return nil
					}

					query = fmt.Sprintf(recValuesPerModule, strings.Join(recordIDs, ","))
					// println(query)
					rows, err = s.DB.QueryContext(ctx, query)
					if err != nil {
						return
					}

					defer func() {
						// assign error to return value...
						err = rows.Close()
					}()

					for rows.Next() {
						if err = rows.Err(); err != nil {
							return
						}

						err = rows.Scan(&recordID, &name, &value, &ref, &place)
						if err != nil {
							return
						}

						sliceLastRecordID = recordID
						if values[recordID] == nil {
							values[recordID] = make(map[string][]any)
						}

						// mimicking behaviour of
						// SimpleJsonDocColumn.Encode function
						field = fields.FindByName(name)
						if field == nil {
							continue
						}

						if !field.Multi && len(values[recordID][name]) > 0 {
							// constraint single-value fields
							continue
						}

						switch {
						case field.IsBoolean():
							intVal = cast.ToBool(value)
						default:
							intVal = value
						}

						values[recordID][name] = append(values[recordID][name], intVal)
						sliceLastRecordID = recordID
					}

					return
				}()

				if err != nil {
					return
				}

				// Update records with collected values
				var encoded []byte
				for ID, kv := range values {
					if len(values) == 0 {
						return nil
					}

					encoded, err = json.Marshal(kv)
					if err != nil {
						return err
					}

					upd := s.Dialect.GOQU().
						Update(model.Record.Ident).
						// postgresql gets a bit confused
						Prepared(false).
						Where(exp.Ex{"id": ID}).
						Set(exp.Record{"values": encoded})

					sql, aa, err := upd.ToSQL()
					if err != nil {
						return err
					}

					_, err = s.DB.ExecContext(ctx, sql, aa...)
					_ = aa
					_ = sql
					_ = upd

					if err != nil {
						return err
					}
				}

				countRecords += len(values)

				perModLog.Debug("migrating record values",
					zap.Int("records", len(values)),
					zap.Duration("dur", time.Now().Sub(bmStart).Round(time.Millisecond)),
					zap.Float64("%", float64(countRecords)/float64(totalRecords)*100),
				)

				if len(values) < recordSliceSize {
					break
				}
			}

			return nil
		}()

		if err != nil {
			return
		}
	}

	err = dropTable(ctx, s, "compose_record_value")
	if err != nil {
		return err
	}

	log.Debug("compose_record_value table removed")
	return nil
}

func fix_2022_09_00_addRevisionOnComposeRecords(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"compose_record",
		model.Record.Attributes.FindByIdent("Revision"),
	)
}

func fix_2022_09_00_addMetaOnComposeRecords(ctx context.Context, s *Store) (err error) {
	var (
		log = s.log(ctx)

		groupedMeta = make(map[uint64]map[string]any)
		packed      []byte
	)

	_, err = s.DataDefiner.TableLookup(ctx, "labels")
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	err = addColumn(ctx, s,
		"compose_record",
		model.Record.Attributes.FindByIdent("Meta"),
	)

	if err != nil {
		return
	}

	return s.Tx(ctx, func(ctx context.Context, s store.Storer) (err error) {
		log.Info("collecting record labels")
		ll, _, err := store.SearchLabels(ctx, s, labelsType.LabelFilter{Kind: "compose:record"})
		if err != nil {
			return
		}

		log.Info("grouping labels", zap.Int("count", len(ll)))
		for _, l := range ll {
			if _, has := groupedMeta[l.ResourceID]; !has {
				groupedMeta[l.ResourceID] = make(map[string]any)
			}

			groupedMeta[l.ResourceID][l.Name] = l.Value
			if err = store.DeleteLabel(ctx, s, l); err != nil {
				return
			}
		}

		log.Info("updating records with meta", zap.Int("count", len(ll)))
		for recordID, labels := range groupedMeta {
			packed, err = json.Marshal(labels)
			_, err = s.(*Store).DB.ExecContext(ctx, "UPDATE compose_record SET meta = $1 WHERE id = $2", packed, recordID)
			if err != nil {
				return
			}
		}

		return

	})
}

func fix_2022_09_00_addMissingNodeIdOnFederationMapping(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"federation_module_mapping",
		&dal.Attribute{Ident: "node_id", Type: &dal.TypeID{}},
	)
}

func fix_2022_09_07_changePostgresIdColumnsDatatype(ctx context.Context, s *Store) (err error) {
	var tableName string
	if !strings.HasPrefix(s.DB.DriverName(), "postgres") {
		return
	}

	log := s.log(ctx)
	log.Info("changing postgres ID columns")

	tnames := tableNames()
	tnamesQry := `SELECT table_name FROM INFORMATION_SCHEMA.COLUMNS  WHERE column_name = 'id' AND
                table_name IN('` + strings.Join(tnames, "', '") + `') AND table_name NOT LIKE 'auth_sessions' AND is_updatable = 'YES'`

	rows, err := s.DB.QueryContext(ctx, tnamesQry)
	if err != nil {
		return err
	}

	for rows.Next() {
		// Get the table name
		err = rows.Scan(&tableName)
		if err != nil {
			return err
		}

		query := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN id TYPE NUMERIC USING CAST(id AS NUMERIC)", tableName)
		_, err = s.DB.ExecContext(ctx, query)
		if err != nil {
			return err
		}
	}

	return nil
}

func fix_2023_03_00_migrateComposeModuleConfigForRecordDeDup(ctx context.Context, s *Store) (err error) {
	// @note skipping this for SQL server since it was introduced with 2023.3.0 so there.
	// There are issues with the implementation which won't work on mssql.
	// Since there is no way this would do anything on mssql, we can skip it.
	if s.DB.DriverName() == "sqlserver" {
		return
	}

	type (
		oldRule struct {
			Name       string   `json:"name"`
			Strict     bool     `json:"strict"`
			Attributes []string `json:"attributes"`
		}
		rules struct {
			Rules []oldRule `json:"rules"`
		}
	)

	var (
		log     = s.log(ctx)
		query   string
		aux     []byte
		rr      rules
		rows    *sql.Rows
		modules types.ModuleSet
	)

	_, err = s.DataDefiner.TableLookup(ctx, model.Module.Ident)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Debug("skipping module config recordDeDup migration: compose_module table not found")
			return nil
		}
		return err
	}

	const (
		moduleConfigRecordDeDup = `
			SELECT compose_module.config -> 'recordDeDup' AS recordDeDup
			FROM compose_module
			WHERE compose_module.id = %d`
	)

	modules, _, err = s.SearchComposeModules(ctx, types.ModuleFilter{})
	if err != nil {
		return
	}

	// 1. Check if module has recordDeDup rules
	// 2. If yes, migrate them to new format
	// 3. Save module
	for _, m := range modules {
		var (
			migratedRules types.DeDupRuleSet
		)

		if err = s.Tx(ctx, func(ctx context.Context, s store.Storer) (err error) {
			log.Debug("collecting module.config.recordDeDup for module", logger.Uint64("id", m.ID))

			query = fmt.Sprintf(moduleConfigRecordDeDup, m.ID)
			rows, err = s.(*Store).DB.QueryContext(ctx, query)
			if err != nil {
				return
			}

			defer func() {
				// assign error to return value...
				err = rows.Close()
			}()

			for rows.Next() {
				if err = rows.Err(); err != nil {
					log.Info("failed to scan rows to migrated module.config.recordDeDup for module",
						logger.Uint64("id", m.ID))
					return
				}

				err = rows.Scan(&aux)
				if err != nil {
					continue
				}

				err = json.Unmarshal(aux, &rr)
				if err != nil {
					continue
				}
			}

			for _, r := range rr.Rules {
				if len(r.Attributes) == 0 {
					continue
				}

				var rcc types.DeDupRuleConstraintSet
				for _, atr := range r.Attributes {
					if len(atr) == 0 {
						continue
					}

					rcc = append(rcc, &types.DeDupRuleConstraint{
						Attribute:  atr,
						Modifier:   "ignore-case",
						MultiValue: "equal",
					})
				}

				migratedRules = append(migratedRules, &types.DeDupRule{
					Strict:        r.Strict,
					ConstraintSet: rcc,
				})
			}

			if len(migratedRules) > 0 {
				m.Config.RecordDeDup.Rules = migratedRules

				log.Debug("saving migrated module.config.recordDeDup for module", logger.Uint64("id", m.ID))
				if err = s.UpdateComposeModule(ctx, m); err != nil {
					log.Debug("error saving migrated module.config.recordDeDup for module", logger.Uint64("id", m.ID))
					return
				}
			}

			return
		}); err != nil {
			continue
		}
	}

	return
}

func fix_2024_09_03_dropFederationNodeSyncPrimaryKey(ctx context.Context, s *Store) (err error) {
	// confirm that the table exists first
	_, err = s.DataDefiner.TableLookup(ctx, "federation_nodes_sync")
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	var (
		alterQry   string
		rows       *sql.Rows
		tableName  = "federation_nodes_sync"
		driverName = s.DB.DriverName()
	)

	switch {
	case strings.HasPrefix(driverName, "sqlite"):
		tempTable := tableName + "_temp"
		sqlStatements := []string{
			fmt.Sprintf("CREATE TABLE %s AS SELECT * FROM %s WHERE 1=0", tempTable, tableName),
			fmt.Sprintf("INSERT INTO %s SELECT * FROM %s", tempTable, tableName),
			fmt.Sprintf("DROP TABLE %s", tableName),
			fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tempTable, tableName),
		}

		if err = s.Tx(ctx, func(ctx context.Context, s store.Storer) (err error) {
			for _, sql := range sqlStatements {
				if _, err = s.(*Store).DB.ExecContext(ctx, sql); err != nil {
					return err
				}
			}

			return
		}); err != nil {
			return err
		}

		return nil
	case strings.HasPrefix(driverName, "mysql"):
		// check if the primary key exists
		pkExistsQry := `SELECT COUNT(*) as pk_exists FROM information_schema.table_constraints
        WHERE table_name = 'federation_nodes_sync' AND constraint_type = 'PRIMARY KEY'`
		rows, err = s.DB.QueryContext(ctx, pkExistsQry)
		if err != nil {
			return err
		}
		defer rows.Close()
		var pkCount int
		if rows.Next() {
			err := rows.Scan(&pkCount)
			if err != nil {
				return err
			}
		}
		// if the primary key doesn't exists, we skip the drop
		if pkCount == 0 {
			return nil
		}
		alterQry = `ALTER TABLE federation_nodes_sync DROP PRIMARY KEY`

	case strings.HasPrefix(driverName, "sqlserver"):
		pkQry := `SELECT CONSTRAINT_NAME AS primary_key FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
		    WHERE CONSTRAINT_TYPE='PRIMARY KEY' AND TABLE_NAME='federation_nodes_sync'`

		rows, err = s.DB.QueryContext(ctx, pkQry)
		if err != nil {
			return err
		}
		defer rows.Close()

		var pkName string
		if rows.Next() {
			err := rows.Scan(&pkName)
			if err != nil {
				return err
			}
		}

		if pkName == "" {
			return nil
		}
		alterQry = fmt.Sprintf("ALTER TABLE federation_nodes_sync DROP CONSTRAINT %s", pkName)

	case strings.HasPrefix(driverName, "postgres"):
		alterQry = `ALTER TABLE federation_nodes_sync DROP CONSTRAINT IF EXISTS federation_nodes_sync_pkey`
	}

	if _, err = s.DB.ExecContext(ctx, alterQry); err != nil {
		return fmt.Errorf("failed to drop primary key constraint on federation_nodes_sync: %w", err)
	}

	return nil
}

func fix_2024_09_03_renameFederationNodeSyncNodeID(ctx context.Context, s *Store) (err error) {
	return renameColumn(ctx, s, "federation_nodes_sync", "node_id", "rel_node")
}

func fix_2024_09_03_renameFederationNodeSyncComposeID(ctx context.Context, s *Store) (err error) {
	return renameColumn(ctx, s, "federation_nodes_sync", "module_id", "rel_module")
}

func fix_2024_09_05_addUserGroupReferenceToUser(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"users",
		&dal.Attribute{
			Ident: "UserGroupID",
			Type: &dal.TypeRef{Nullable: true, HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user-group",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_user_group"},
		},
	)
}

func fix_2024_09_05_addRelResourceRoleMembershipColumn(ctx context.Context, s *Store) (err error) {
	err, exists := addColumnExists(ctx, s,
		"role_members",
		&dal.Attribute{
			Ident: "rel_resource",
			Type:  &dal.TypeText{Nullable: true, HasDefault: false},
			Store: &dal.CodecAlias{Ident: "rel_resource"},
		},
	)
	if err != nil {
		return
	}

	if exists {
		return
	}

	db := s.DB.(goqu.SQLDatabase)
	sql, params, err := s.Dialect.
		GOQU().
		DB(db).
		Update("role_members").
		Set(goqu.Record{
			"rel_resource": goqu.Func("concat",
				goqu.L("'corteza::system:user/'"),
				goqu.C("rel_user"),
			),
		}).
		ToSQL()

	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, params...)
	if err != nil {
		return err
	}

	return
}

func fix_2024_09_05_addRelResourceColumn(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"role_members",
		&dal.Attribute{
			Ident: "rel_resource",
			Type:  &dal.TypeText{Nullable: true},
			Store: &dal.CodecAlias{Ident: "rel_resource"},
		},
	)
}

func fix_2024_09_03_addFederationNodeSyncNodeIDIndex(ctx context.Context, s *Store) (err error) {
	var (
		tableName  = "federation_nodes_sync"
		indexName  = "federation_nodes_sync_idxRelNode"
		driverName = s.DB.DriverName()
	)

	// confirm that the table exists
	_, err = s.DataDefiner.TableLookup(ctx, tableName)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if strings.HasPrefix(driverName, "mysql") || strings.HasPrefix(driverName, "sqlserver") {
		indx, err := s.DataDefiner.IndexLookup(ctx, indexName, tableName)
		if err != nil {
			return err
		}

		// if the index already exists, skip the creation
		if indx != nil {
			return nil
		}
	}

	nodeIDIndx := ddl.Index{
		TableIdent: tableName,
		Ident:      indexName,
		Type:       "BTREE",

		Fields: []*ddl.IndexField{
			{
				Column: "rel_node",
			},
		},
	}

	return s.DataDefiner.IndexCreate(ctx, "federation_nodes_sync", &nodeIDIndx)
}
func fix_2024_9_7_migrateLabelsValueToJsonbPostgres(ctx context.Context, s *Store) (err error) {
	var (
		log        = s.log(ctx)
		columnType string
		rows       *sql.Rows
	)
	checkQuery := `
  			SELECT data_type
  			FROM information_schema.columns
  			WHERE table_name = 'labels' AND column_name = 'value'`

	exists, err := func() (bool, error) {
		rows, err = s.DB.QueryContext(ctx, checkQuery)
		if err != nil {
			return false, err
		}
		defer rows.Close()

		if rows.Next() {
			if err = rows.Scan(&columnType); err != nil {
				return false, err
			}
		}

		if len(columnType) == 0 {
			log.Debug("labels table not found, skipping migration")
			return true, nil
		}

		if columnType == "jsonb" {
			log.Debug("labels.value column already jsonb, skipping migration")
			return true, nil
		}

		return false, nil
	}()

	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	log.Info("migrating labels.value column from text to jsonb")

	updateQuery := `
  			UPDATE labels
  			SET value = jsonb_build_object('value', value::text)
  			WHERE value IS NOT NULL
  			  AND value::text NOT LIKE '{%'`

	if _, err = s.DB.ExecContext(ctx, updateQuery); err != nil {
		return err
	}

	alterQuery := `ALTER TABLE labels ALTER COLUMN value TYPE jsonb USING value::jsonb`
	if _, err = s.DB.ExecContext(ctx, alterQuery); err != nil {
		return err
	}

	log.Info("successfully migrated labels.value column to jsonb")
	return nil

}
func fix_2024_9_7_migrateLabelsValueToJsonbMySql(ctx context.Context, s *Store) (err error) {
	var (
		log        = s.log(ctx)
		columnType string
		rows       *sql.Rows
	)
	checkQuery := `
  			SELECT data_type
  			FROM information_schema.columns
  			WHERE table_name = 'labels' AND column_name = 'value'`

	exists, err := func() (bool, error) {
		rows, err = s.DB.QueryContext(ctx, checkQuery)
		if err != nil {
			return false, err
		}
		defer rows.Close()

		if rows.Next() {
			if err = rows.Scan(&columnType); err != nil {
				return false, err
			}
		}

		if len(columnType) == 0 {
			log.Debug("labels table not found, skipping migration")
			return true, nil
		}

		if columnType == "json" {
			log.Debug("labels.value column already json, skipping migration")
			return true, nil
		}

		return false, nil
	}()

	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	log.Info("migrating labels.value column from text to json")

	updateQuery := `
  			UPDATE labels
  			SET value = JSON_OBJECT('value', value)
  			WHERE value IS NOT NULL
  			  AND value NOT LIKE '{%'`

	if _, err = s.DB.ExecContext(ctx, updateQuery); err != nil {
		return err
	}

	alterQuery := `ALTER TABLE labels MODIFY COLUMN value JSON`
	if _, err = s.DB.ExecContext(ctx, alterQuery); err != nil {
		return err
	}

	log.Info("successfully migrated labels.value column to json")
	return nil
}
func fix_2024_9_7_migrateLabelsValueToJsonbSqlite(ctx context.Context, s *Store) (err error) {
	var (
		log        = s.log(ctx)
		columnType string
		rows       *sql.Rows
	)

	checkQuery := `
  			SELECT type
  			FROM pragma_table_info('labels')
  			WHERE name = 'value'`

	exists, err := func() (bool, error) {
		rows, err = s.DB.QueryContext(ctx, checkQuery)
		if err != nil {
			return false, err
		}
		defer rows.Close()

		if rows.Next() {
			if err = rows.Scan(&columnType); err != nil {
				return false, err
			}
		}

		if len(columnType) == 0 {
			log.Debug("labels table not found, skipping migration")
			return true, nil
		}

		if strings.ToUpper(columnType) == "JSON" {
			log.Debug("labels.value column already json, skipping migration")
			return true, nil
		}

		return false, nil
	}()

	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	log.Info("migrating labels.value data format for sqlite")

	updateQuery := `
  			UPDATE labels
  			SET value = json_object('value', value)
  			WHERE value IS NOT NULL
  			  AND value NOT LIKE '{%'`

	if _, err = s.DB.ExecContext(ctx, updateQuery); err != nil {
		return err
	}

	tempTable := "labels_temp"
	sqlStatements := []string{
		fmt.Sprintf("CREATE TABLE %s (kind TEXT, rel_resource NUMERIC, name TEXT, value JSON, PRIMARY KEY (kind, rel_resource, name))", tempTable),
		fmt.Sprintf("INSERT INTO %s SELECT * FROM labels", tempTable),
		"DROP TABLE labels",
		fmt.Sprintf("ALTER TABLE %s RENAME TO labels", tempTable),
	}

	if err = s.Tx(ctx, func(ctx context.Context, s store.Storer) (err error) {
		for _, sql := range sqlStatements {
			if _, err = s.(*Store).DB.ExecContext(ctx, sql); err != nil {
				return err
			}
		}
		return
	}); err != nil {
		return err
	}

	log.Info("successfully migrated labels.value data format")
	return nil
}
func fix_2024_9_7_migrateLabelsValueToJsonbSqlserver(ctx context.Context, s *Store) (err error) {
	var (
		log        = s.log(ctx)
		columnType string
		rows       *sql.Rows
	)
	checkQuery := `
  			SELECT DATA_TYPE
  			FROM INFORMATION_SCHEMA.COLUMNS
  			WHERE TABLE_NAME = 'labels' AND COLUMN_NAME = 'value'`

	exists, err := func() (bool, error) {
		rows, err = s.DB.QueryContext(ctx, checkQuery)
		if err != nil {
			return false, err
		}
		defer rows.Close()

		if rows.Next() {
			if err = rows.Scan(&columnType); err != nil {
				return false, err
			}
		}

		if len(columnType) == 0 {
			log.Debug("labels table not found, skipping migration")
			return true, nil
		}

		if columnType == "nvarchar" {
			log.Debug("labels.value column already nvarchar, skipping migration")
			return true, nil
		}

		return false, nil
	}()

	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	log.Info("migrating labels.value data format for sqlserver")

	updateQuery := `
  			UPDATE labels
  			SET value = JSON_OBJECT('value', value)
  			WHERE value IS NOT NULL
  			  AND value NOT LIKE '{%'`

	if _, err = s.DB.ExecContext(ctx, updateQuery); err != nil {
		return err
	}

	alterQuery := `ALTER TABLE labels ALTER COLUMN value NVARCHAR(MAX)`
	if _, err = s.DB.ExecContext(ctx, alterQuery); err != nil {
		return err
	}

	log.Info("successfully migrated labels.value data format")
	return nil
}
func fix_2024_9_7_migrateLabelsValueToJsonb(ctx context.Context, s *Store) (err error) {
	var driverName = s.DB.DriverName()

	switch {
	case strings.HasPrefix(driverName, "postgres"):
		return fix_2024_9_7_migrateLabelsValueToJsonbPostgres(ctx, s)

	case strings.HasPrefix(driverName, "mysql"):
		return fix_2024_9_7_migrateLabelsValueToJsonbMySql(ctx, s)

	case strings.HasPrefix(driverName, "sqlite"):
		return fix_2024_9_7_migrateLabelsValueToJsonbSqlite(ctx, s)

	case strings.HasPrefix(driverName, "sqlserver"):
		return fix_2024_9_7_migrateLabelsValueToJsonbSqlserver(ctx, s)
	}
	return nil

}
func fix_2026_07_00_addRevisionColumnsOnProjects(ctx context.Context, s *Store) error {
	if err := addColumn(ctx, s, "projects", &dal.Attribute{
		Ident: "ProjectID",
		Type:  &dal.TypeID{HasDefault: true, DefaultValue: 0},
		Store: &dal.CodecAlias{Ident: "root_project_id"},
	}); err != nil {
		return err
	}
	if err := addColumn(ctx, s, "projects", &dal.Attribute{
		Ident: "ParentRevisionID",
		Type:  &dal.TypeID{HasDefault: true, DefaultValue: 0},
		Store: &dal.CodecAlias{Ident: "parent_revision_id"},
	}); err != nil {
		return err
	}
	return addColumn(ctx, s, "projects", &dal.Attribute{
		Ident: "Revision",
		Type:  &dal.TypeNumber{HasDefault: true, DefaultValue: 0, Precision: 0},
		Store: &dal.CodecAlias{Ident: "revision"},
	})
}

// fix_2026_07_14_addModeOnProjects promotes project mode out of the config JSON
// to a top-level sortable column. The column defaults to 'free'; existing rows
// are backfilled from the old config->>'mode' so gated projects keep their mode.
func fix_2026_07_14_addModeOnProjects(ctx context.Context, s *Store) error {
	if err := addColumn(ctx, s, "projects", &dal.Attribute{
		Ident: "Mode",
		Type:  &dal.TypeText{Length: 32, HasDefault: true, DefaultValue: "free"},
		Store: &dal.CodecAlias{Ident: "mode"},
	}); err != nil {
		return err
	}

	// Skip the backfill on a fresh DB (table not created yet — the column comes
	// from the model, and there are no rows to migrate).
	if _, err := s.DataDefiner.TableLookup(ctx, "projects"); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Dialect-aware extraction of config.mode from the JSON column.
	var expr string
	switch {
	case strings.HasPrefix(s.DB.DriverName(), "postgres"):
		expr = "config->>'mode'"
	case strings.HasPrefix(s.DB.DriverName(), "mysql"):
		expr = "JSON_UNQUOTE(JSON_EXTRACT(config, '$.mode'))"
	case strings.HasPrefix(s.DB.DriverName(), "sqlite"):
		expr = "json_extract(config, '$.mode')"
	case s.DB.DriverName() == "sqlserver":
		expr = "JSON_VALUE(config, '$.mode')"
	default:
		return nil
	}

	q := fmt.Sprintf("UPDATE projects SET mode = %s WHERE %s IS NOT NULL AND %s <> ''", expr, expr, expr)
	_, err := s.DB.ExecContext(ctx, q)
	return err
}

// fix_2026_07_21_dropModeOnProjects removes the projects.mode column. Build
// modes (free/gated) were dropped from the product — all projects now behave
// identically and go through the same publish-approval flow — so the column is
// obsolete. Registered after fix_2026_07_14_addModeOnProjects so upgrades from
// any prior version add-then-drop cleanly; dropColumns no-ops when the column
// is already absent (fresh DBs, where the model never defined it).
func fix_2026_07_21_dropModeOnProjects(ctx context.Context, s *Store) error {
	return dropColumns(ctx, s, "projects", "mode")
}

// fix_2026_07_28_addRelRevisionOnProjectWorkItems adds the rel_revision column
// (+ BTREE index) to the six project work-item tables: project backlog items
// and the five category tables (incident/feature/task/privacy/review). Work
// items always file against the chain root project (rel_project) and
// optionally point at the revision they're assigned to via rel_revision, so a
// revision reads like a GitHub milestone over one shared item pool. Existing
// rows are left unassigned (null/0) by design — no backfill.
func fix_2026_07_28_addRelRevisionOnProjectWorkItems(ctx context.Context, s *Store) error {
	workItemModels := []struct {
		table string
		model *dal.Model
	}{
		{"project_backlog_items", systemModel.ProjectBacklogItem},
		{"project_incidents", systemModel.ProjectIncident},
		{"project_features", systemModel.ProjectFeature},
		{"project_privacys", systemModel.ProjectPrivacy},
		{"project_reviews", systemModel.ProjectReview},
		{"project_tasks", systemModel.ProjectTask},
	}

	for _, wi := range workItemModels {
		if _, err := s.DataDefiner.TableLookup(ctx, wi.table); err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			return err
		}

		attr := wi.model.Attributes.FindByIdent("RevisionID")
		if attr == nil {
			return fmt.Errorf("%s model is missing the RevisionID attribute", wi.table)
		}
		if err := addColumn(ctx, s, wi.table, attr); err != nil {
			return err
		}

		indexName := wi.table + "_rel_revision"
		idx, err := s.DataDefiner.IndexLookup(ctx, indexName, wi.table)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		if idx != nil {
			continue
		}

		if err := s.DataDefiner.IndexCreate(ctx, wi.table, &ddl.Index{
			TableIdent: wi.table,
			Ident:      indexName,
			Type:       "BTREE",
			Fields:     []*ddl.IndexField{{Column: "rel_revision"}},
		}); err != nil {
			return err
		}
	}

	return nil
}

// fix_2026_07_30_dropGlobalUniqueHandleOnAgents drops the agents_uniqueHandle
// index, superseded by agents_uniqueHandlePerProject (see system/agent.cue).
//
// A project revision branch copies the parent's agents into the draft (see
// system/service/project_revision_clone.go), so both revisions necessarily hold
// an agent with the same handle at once -- impossible while the handle is
// globally unique. Carrying the handle over is not cosmetic: the publish diff
// identifies resources across revisions by kind + handle (diffSources), so a
// renamed copy would read as a removal plus an addition.
//
// This fix is required because the upgrade only ever ADDS indexes the model
// declares; it never drops ones the model dropped. Without it an already
// -migrated database keeps enforcing the stale global index and every branch
// fails with "not unique", while a fresh database works -- the nastiest kind of
// environment-dependent bug.
func fix_2026_07_30_dropGlobalUniqueHandleOnAgents(ctx context.Context, s *Store) error {
	return dropIndexes(ctx, s, "agents", "agents_uniqueHandle")
}

// fix_2026_07_30_dropGlobalUniqueHandleOnChatbots is the agent fix above,
// applied to chatbots for the same reason: the revision branch now copies a
// project's chatbots too (see system/service/project_revision_clone.go), so
// parent and draft hold a same-handled chatbot at once.
//
// The handle has to be carried over rather than suffixed because the publish
// diff matches resources across revisions on kind + handle; a renamed copy
// would read as "chatbot removed, chatbot added" on a branch that changed
// nothing. chatbots_uniqueHandlePerProject (system/chatbot.cue) replaces it.
//
// Separate from the model change because the upgrade only ever ADDS indexes the
// model declares — it never drops the ones the model dropped. Without this, a
// database migrated before today keeps enforcing the stale global index and
// every branch of a project with a chatbot fails with "not unique", while a
// freshly created one works.
func fix_2026_07_30_dropGlobalUniqueHandleOnChatbots(ctx context.Context, s *Store) error {
	return dropIndexes(ctx, s, "chatbots", "chatbots_uniqueHandle")
}

// fix_2026_07_30_addArchivedAtOnProjects splits archiving out of the project
// status. Archiving is the one lifecycle change a user makes directly, and
// while it shared a column with the publish machinery's status, the generic
// update had to accept a status from the client -- which is how a live project
// could be flipped back to "draft" and then have its schema edited underneath
// its own records.
//
// Rows archived under the old scheme are carried over: without that they would
// come back off the shelf on upgrade -- every listing filters on archived_at
// now -- while still carrying a status the project resolver refuses, so they
// would be visible and un-openable at the same time.
//
// The status they had BEFORE being archived is not recoverable; it was
// overwritten when they were archived, which is the flaw that motivated the
// split. They come back as drafts, which is the safe reading: a draft is
// editable and not live, and re-publishing is a deliberate act.
func fix_2026_07_30_addArchivedAtOnProjects(ctx context.Context, s *Store) error {
	if err := addColumn(ctx, s, "projects", &dal.Attribute{
		Ident: "ArchivedAt",
		Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
		Store: &dal.CodecAlias{Ident: "archived_at"},
	}); err != nil {
		return err
	}

	if _, err := s.DataDefiner.TableLookup(ctx, "projects"); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	_, err := s.DB.ExecContext(ctx, `
		UPDATE projects
		   SET archived_at = COALESCE(updated_at, created_at),
		       status = 'draft'
		 WHERE status = 'archived' AND archived_at IS NULL`)

	return err
}

// fix_2026_07_31_addApprovalOnProjects gives the publish approval cycle a home
// on the revision row.
//
// It used to live in a Pinia ref in the browser: it did not survive a reload, a
// second user never saw a submitted request, a submitter could approve their
// own work, and the server published on confirm=true alone -- so the entire
// gate was advisory. Publishing is now refused unless the row says approved.
//
// Existing rows come back as drafts (the zero value of the new column), which
// is the safe reading: nothing that was never reviewed should be publishable
// because it predates the review. Already-live revisions are unaffected --
// publish only ever looks at a draft.
func fix_2026_07_31_addApprovalOnProjects(ctx context.Context, s *Store) error {
	// The three text columns carry a DDL default even though the model does
	// not: they are NOT NULL, and postgres refuses to add a NOT NULL column to
	// a table that already has rows unless it is told what those rows should
	// say. 'draft' is the honest answer for every existing project -- nothing
	// that predates the gate was ever reviewed.
	cols := []*dal.Attribute{{
		Ident: "ApprovalStatus",
		Type:  &dal.TypeText{Length: 32, HasDefault: true, DefaultValue: "draft"},
		Store: &dal.CodecAlias{Ident: "approval_status"},
	}, {
		Ident: "ApprovalPlan",
		Type:  &dal.TypeText{Length: 64, HasDefault: true, DefaultValue: ""},
		Store: &dal.CodecAlias{Ident: "approval_plan"},
	}, {
		Ident: "ApprovalNote",
		Type:  &dal.TypeText{HasDefault: true, DefaultValue: ""},
		Store: &dal.CodecAlias{Ident: "approval_note"},
	}, {
		Ident: "ApprovalSubmittedBy",
		Type:  &dal.TypeID{HasDefault: true, DefaultValue: 0},
		Store: &dal.CodecAlias{Ident: "approval_submitted_by"},
	}, {
		Ident: "ApprovalSubmittedAt",
		Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
		Store: &dal.CodecAlias{Ident: "approval_submitted_at"},
	}, {
		Ident: "ApprovalDecidedBy",
		Type:  &dal.TypeID{HasDefault: true, DefaultValue: 0},
		Store: &dal.CodecAlias{Ident: "approval_decided_by"},
	}, {
		Ident: "ApprovalDecidedAt",
		Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
		Store: &dal.CodecAlias{Ident: "approval_decided_at"},
	}}

	// Skipped wholesale on a store that has not created the table yet -- the
	// same guard the archived_at fix carries, for the same reason: a fresh
	// install builds the table from the model, which already has the columns.
	if _, err := s.DataDefiner.TableLookup(ctx, "projects"); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	for _, col := range cols {
		if err := addColumn(ctx, s, "projects", col); err != nil {
			return err
		}
	}

	return nil
}

// fix_2026_07_30_backfillProjectRefOnComposeResources aligns every compose
// resource with the project its namespace belongs to.
//
// Only charts derived rel_project from their namespace; modules and pages took
// it from the request payload, so anything created by a client that did not
// send one -- every module and page made through the API directly, and every
// resource an envoy clone copied -- carries 0. Nothing that enumerates a
// project by rel_project could see them: the resource graph showed a project
// with modules as having none and reported every reference into them as
// dangling, and the deployment plan built on that graph could not report a
// page as added or removed at all.
//
// The namespace is the authority: a compose resource is addressed through it,
// so its project is the namespace's project, full stop.
func fix_2026_07_30_backfillProjectRefOnComposeResources(ctx context.Context, s *Store) error {
	// Most of these carry rel_namespace and can read the project straight off
	// it; module fields hold only rel_module, so they take the same route one
	// hop later. Either way the project comes from the namespace, never from
	// a sibling row that might itself still be an orphan.
	sources := map[string]string{
		"compose_module":      "SELECT ns.rel_project FROM compose_namespace AS ns WHERE ns.id = compose_module.rel_namespace",
		"compose_page":        "SELECT ns.rel_project FROM compose_namespace AS ns WHERE ns.id = compose_page.rel_namespace",
		"compose_page_layout": "SELECT ns.rel_project FROM compose_namespace AS ns WHERE ns.id = compose_page_layout.rel_namespace",
		"compose_chart":       "SELECT ns.rel_project FROM compose_namespace AS ns WHERE ns.id = compose_chart.rel_namespace",
		"compose_module_field": "SELECT ns.rel_project FROM compose_namespace AS ns" +
			" JOIN compose_module AS m ON m.rel_namespace = ns.id" +
			" WHERE m.id = compose_module_field.rel_module",
	}

	// The WHERE clause is what makes this both complete and idempotent, and an
	// earlier version got both wrong. "rel_project = 0" missed the rows that
	// matter most -- a resource copied by an envoy clone carries the PARENT
	// project's id, not 0, so every revision branched before the repoint
	// existed stayed mis-attributed -- while also matching every legitimately
	// project-less resource on a classic install and writing 0 back over 0 on
	// every single boot.
	//
	// So: adopt where the namespace says a project and the row disagrees, and
	// only there. A namespace with no project leaves its resources alone.
	predicates := map[string]string{
		"compose_module":      "compose_module.rel_namespace",
		"compose_page":        "compose_page.rel_namespace",
		"compose_page_layout": "compose_page_layout.rel_namespace",
		"compose_chart":       "compose_chart.rel_namespace",
	}

	// Map iteration is unordered and this writes to the database; fix an order
	// so a failure reports the same table every time.
	tables := make([]string, 0, len(sources))
	for table := range sources {
		tables = append(tables, table)
	}
	sort.Strings(tables)

	for _, table := range tables {
		if _, err := s.DataDefiner.TableLookup(ctx, table); err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			return err
		}

		var owner string
		if nsCol, ok := predicates[table]; ok {
			owner = fmt.Sprintf(
				"SELECT 1 FROM compose_namespace AS ns WHERE ns.id = %s"+
					" AND ns.rel_project > 0 AND ns.rel_project <> %s.rel_project",
				nsCol, table,
			)
		} else {
			owner = "SELECT 1 FROM compose_namespace AS ns" +
				" JOIN compose_module AS m ON m.rel_namespace = ns.id" +
				" WHERE m.id = compose_module_field.rel_module" +
				" AND ns.rel_project > 0 AND ns.rel_project <> compose_module_field.rel_project"
		}

		q := fmt.Sprintf(
			"UPDATE %s SET rel_project = (%s) WHERE EXISTS (%s)",
			table, sources[table], owner,
		)

		if _, err := s.DB.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("backfill rel_project on %s: %w", table, err)
		}
	}

	return nil
}

func count(ctx context.Context, s *Store, table string, ee ...goqu.Expression) (count int) {
	db := s.DB.(goqu.SQLDatabase)

	_, err := s.Dialect.
		GOQU().
		DB(db).
		Select(goqu.COUNT(goqu.Star())).
		From(table).
		Where(ee...).
		ScanValContext(ctx, &count)

	if err != nil {
		panic(err)
	}

	return
}
