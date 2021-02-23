package rdbms

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/store/rdbms/ddl"
	"go.uber.org/zap"
)

type (
	genericUpgrades struct {
		log *zap.Logger
		u   upgrader
	}

	upgrader interface {
		TableExists(context.Context, string) (bool, error)
		AddColumn(context.Context, string, *ddl.Column) (bool, error)
		DropTable(context.Context, string) (bool, error)
		DropColumn(context.Context, string, string) (bool, error)
		RenameColumn(context.Context, string, string, string) (bool, error)
		AddPrimaryKey(context.Context, string, *ddl.Index) (bool, error)
		Exec(context.Context, string, ...interface{}) error
	}
)

func GenericUpgrades(log *zap.Logger, u upgrader) *genericUpgrades {
	return &genericUpgrades{log, u}
}

func (g genericUpgrades) Before(ctx context.Context) error {
	return g.all(ctx,
		g.RenameActionlog,
		g.RenameReminders,
		g.RenameUsers,
		g.RenameRoles,
		g.RenameRoleMembers,
		g.RenameCredentials,
		g.RenameApplications,
		g.DropOrganisationTable,
	)
}

func (g genericUpgrades) After(ctx context.Context) error {
	return nil
}

func (g genericUpgrades) Upgrade(ctx context.Context, t *ddl.Table) error {
	switch t.Name {
	case "settings":
		return g.all(ctx,
			g.MergeSettingsTables,
		)
	case "rbac_rules":
		return g.all(ctx,
			g.MergePermissionRulesTables,
		)
	case "actionlog":
		return g.all(ctx,
			g.AlterActionlogAddID,
		)
	case "users":
		return g.all(ctx,
			g.AlterUsersDropOrganisation,
			g.AlterUsersDropRelatedUser,
		)
	case "compose_module":
		return g.all(ctx,
			g.AlterComposeModuleRenameJsonToMeta,
		)
	case "compose_module_field":
		return g.all(ctx,
			g.AlterComposeModuleFieldAddExpresions,
		)
	case "messaging_channel":
		return g.all(ctx,
			g.AlterMessagingChannelsDropOrganisation,
		)
	case "messaging_attachment":
		return g.all(ctx,
			g.AlterMessageAttachmentsRenameOwner,
		)
		//case "compose_attachment_binds":
		//	return g.all(ctx,
		//		g.MigrateComposeAttachmentsToBindsTable,
		//	)
	}

	return nil
}

func (genericUpgrades) all(ctx context.Context, ffn ...func(context.Context) error) (err error) {
	for _, fn := range ffn {
		if err = fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

// MergeSettingsTables merges "*_settings" tables into one single "settings"
func (g genericUpgrades) MergeSettingsTables(ctx context.Context) error {
	var (
		err error
		tt  = []struct {
			tbl             string
			applyNamePrefix string
		}{
			{tbl: "sys_settings", applyNamePrefix: ""},
			{tbl: "compose_settings", applyNamePrefix: "compose."},
			{tbl: "messaging_settings", applyNamePrefix: "messaging."},
		}

		// CONCAT does not work in sqlite but we'Ll ignore this since sqlite should
		// not even get execute this query (*_settings tables do not exist)
		merge = `
			INSERT INTO settings (rel_owner, name, value, updated_by, updated_at) 
			SELECT rel_owner, CONCAT('%s', name), value, updated_by, updated_at 
			  FROM %s`
	)

	for _, t := range tt {
		if exists, err := g.u.TableExists(ctx, t.tbl); err != nil {
			return err
		} else if !exists {
			g.log.Debug(fmt.Sprintf("skipping settings merge, table %s already removed", t.tbl))
			continue
		}

		err = g.u.Exec(ctx, fmt.Sprintf(merge, t.applyNamePrefix, t.tbl))
		if err != nil {
			return fmt.Errorf("could not merge %s: %w", t.tbl, err)
		}

		_, err = g.u.DropTable(ctx, t.tbl)
		if err != nil {
			return fmt.Errorf("could not drop %s: %w", t.tbl, err)
		}

		g.log.Debug(fmt.Sprintf("table %s merged into settings and removed", t.tbl))
	}

	return nil
}

// MergeSettingsTables merges "*_settings" tables into one single "settings"
func (g genericUpgrades) MergePermissionRulesTables(ctx context.Context) error {
	var (
		err error
		tt  = []struct {
			tbl string
		}{
			{tbl: "sys_permission_rules"},
			{tbl: "compose_permission_rules"},
			{tbl: "messaging_permission_rules"},
		}

		// CONCAT does not work in sqlite but we'Ll ignore this since sqlite should
		// not even get execute this query (*_permissions tables do not exist)
		merge = `
			INSERT INTO rbac_rules (rel_role, resource, operation, access) 
			SELECT rel_role, resource, operation, access
			  FROM %s`
	)

	for _, t := range tt {
		if exists, err := g.u.TableExists(ctx, t.tbl); err != nil {
			return err
		} else if !exists {
			g.log.Debug(fmt.Sprintf("skipping rbac rules merge, table %s already Deleted", t.tbl))
			continue
		}

		err = g.u.Exec(ctx, fmt.Sprintf(merge, t.tbl))
		if err != nil {
			return fmt.Errorf("could not merge %s: %w", t.tbl, err)
		}

		_, err = g.u.DropTable(ctx, t.tbl)
		if err != nil {
			return fmt.Errorf("could not drop %s: %w", t.tbl, err)
		}

		g.log.Debug(fmt.Sprintf("table %s merged into rbac_rules and Deleted", t.tbl))
	}

	return nil
}

func (g genericUpgrades) RenameActionlog(ctx context.Context) error {
	return g.RenameTable(ctx, "sys_actionlog", "actionlog")
}

func (g genericUpgrades) RenameUsers(ctx context.Context) error {
	return g.RenameTable(ctx, "sys_user", "users")
}

func (g genericUpgrades) RenameRoles(ctx context.Context) error {
	return g.RenameTable(ctx, "sys_role", "roles")
}

func (g genericUpgrades) RenameRoleMembers(ctx context.Context) error {
	return g.RenameTable(ctx, "sys_role_member", "role_members")
}

func (g genericUpgrades) RenameApplications(ctx context.Context) error {
	return g.RenameTable(ctx, "sys_application", "applications")
}

func (g genericUpgrades) RenameCredentials(ctx context.Context) error {
	return g.RenameTable(ctx, "sys_credentials", "credentials")
}

// AlterActionlogAddID adds ID column, fills it with values and adds PK on it
//
// This is MySQL only; other implementations were never in state with actionlog table
// without ID column.
func (g genericUpgrades) AlterActionlogAddID(ctx context.Context) (err error) {
	var (
		added bool
		col   = &ddl.Column{Name: "id", Type: ddl.ColumnType{Type: ddl.ColumnTypeIdentifier}, IsNull: false}
		ind   = &ddl.Index{Fields: []*ddl.IField{{Field: col.Name}}}
		upd   = `UPDATE actionlog SET id = @v := COALESCE(@v, 0) + 1 WHERE id = 0`
	)
	if added, err = g.u.AddColumn(ctx, "actionlog", col); err != nil {
		return err
	} else if !added {
		// not added, no need to continue
		return
	}

	// Now prefill with generated IDs in any case -- if col was added or not
	g.log.Debug(fmt.Sprintf("prefilling missing values for actionlog ID field, might take a while"))
	if err = g.u.Exec(ctx, upd); err != nil {
		return fmt.Errorf("could not prefill actionlog ID field: %w", err)
	}

	if added, err = g.u.AddPrimaryKey(ctx, "actionlog", ind); err != nil {
		return err
	}

	return nil
}

func (g genericUpgrades) RenameReminders(ctx context.Context) error {
	return g.RenameTable(ctx, "sys_reminder", "reminders")
}

func (g genericUpgrades) DropOrganisationTable(ctx context.Context) error {
	_, err := g.u.DropTable(ctx, "organization")
	return err
}

func (g genericUpgrades) AlterUsersDropOrganisation(ctx context.Context) error {
	_, err := g.u.DropColumn(ctx, "users", "rel_organisation")
	return err
}

func (g genericUpgrades) AlterMessagingChannelsDropOrganisation(ctx context.Context) error {
	_, err := g.u.DropColumn(ctx, "messaging_channel", "rel_organisation")
	return err
}

func (g genericUpgrades) AlterUsersDropRelatedUser(ctx context.Context) error {
	_, err := g.u.DropColumn(ctx, "users", "rel_user_id")
	return err
}

func (g genericUpgrades) RenameTable(ctx context.Context, old, new string) error {
	if exists, err := g.u.TableExists(ctx, old); err != nil {
		return err
	} else if !exists {
		g.log.Debug(fmt.Sprintf("skipping %s table rename, old table does not exist", old))
		return nil
	}

	if exists, err := g.u.TableExists(ctx, new); err != nil {
		return err
	} else if exists {
		g.log.Debug(fmt.Sprintf("skipping %s table rename, new table already exist", new))
		return nil
	}

	if err := g.u.Exec(ctx, fmt.Sprintf("ALTER TABLE %s RENAME TO %s", old, new)); err != nil {
		return err
	}

	g.log.Debug(fmt.Sprintf("table %s renamed to %s", old, new))

	return nil
}

//func (g genericUpgrades) MigrateComposeAttachmentsToLinksTable(ctx context.Context) error {
//	var (
//		err error
//		tt  = []struct {
//			tbl string
//		}{
//			{tbl: "sys_permission_rules"},
//			{tbl: "compose_permission_rules"},
//			{tbl: "messaging_permission_rules"},
//		}
//
//		// Are there entries in the attachment_binds table?
//		check = `SELECT COUNT(*) > 0 FROM compose_attachment_links LIMIT 1`
//
//		splitRecordAttachments = `
//			INSERT INTO compose_attachment_links (
//				   rel_namespace, rel_attachment, kind,
//                   ref,
//				   owned_by
//				   created_by, updated_by, deleted_by,
//				   created_at, updated_at, deleted_at
//			)
//			SELECT rel_namespace, rel_attachment, kind,
//				   CASE WHEN kind = 'page'   THEN
//(SELECT
//                        WHEN kind = 'record' THEN 2
//                        ELSE 0 END,
//				   owned_by, 0, 0,
//				   created_at, updated_at, deleted_at
//			  FROM compose_attachment
//                   INNER JOIN compose_record_Value`
//
//		splitPageAttachments = `
//			INSERT INTO compose_attachment_links (
//				   rel_namespace, rel_attachment, kind,
//                   ref,
//				   owned_by
//				   created_by, updated_by, deleted_by,
//				   created_at, updated_at, deleted_at
//			)
//			SELECT rel_namespace, rel_attachment, kind,
//				   CASE WHEN kind = 'page'   THEN
//(SELECT
//                        WHEN kind = 'record' THEN 2
//                        ELSE 0 END,
//				   owned_by, 0, 0,
//				   created_at, updated_at, deleted_at
//			  FROM compose_attachment`
//	)
//
//	g.log.Debug("splitting parts of compose_attachment to compose_attachment_links")
//	err = g.u.Exec(ctx, split)
//	if err != nil {
//		return fmt.Errorf("failed to split compose_attachment: %w", err)
//	}
//
//	for _, col := range []string{"rel_namespace", "kind"} {
//		_, err = g.u.DropColumn(ctx, "compose_attachment", col)
//		if err != nil {
//			return fmt.Errorf("could not drop column compose_attachment.%s: %w", col, err)
//		}
//	}
//
//	g.log.Debug("compose_attachment split to compose_attachment_links")
//
//	return nil
//}

func (g genericUpgrades) AlterComposeModuleRenameJsonToMeta(ctx context.Context) error {
	_, err := g.u.RenameColumn(ctx, "compose_module", "json", "meta")
	return err
}

func (g genericUpgrades) AlterComposeModuleFieldAddExpresions(ctx context.Context) (err error) {
	var (
		col = &ddl.Column{
			Name:         "expressions",
			Type:         ddl.ColumnType{Type: ddl.ColumnTypeJson},
			IsNull:       false,
			DefaultValue: "'{}'",
		}
	)

	_, err = g.u.AddColumn(ctx, "compose_module_field", col)
	return
}

func (g genericUpgrades) AlterMessageAttachmentsRenameOwner(ctx context.Context) error {
	_, err := g.u.RenameColumn(ctx, "messaging_attachment", "rel_user", "rel_owner")
	return err
}
