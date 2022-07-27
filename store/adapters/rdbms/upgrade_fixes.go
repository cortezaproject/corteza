package rdbms

import (
	"context"
	. "github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
)

func fix202209_extendComposeModuleForPrivacyAndDAL(ctx context.Context, s *Store) (err error) {
	s.log(ctx).Info("extending compose_module table with config column")
	return s.SchemaAPI.AddColumn(
		ctx, s.DB,
		&Table{Name: "compose_module"},
		&Column{Type: ColumnType{Type: ColumnTypeJson}, DefaultValue: "'{}'", Name: "config"},
	)
}

func fix202209_extendComposeModuleFieldsForPrivacyAndDAL(ctx context.Context, s *Store) (err error) {
	s.log(ctx).Info("extending compose_module_field table with config column")
	return s.SchemaAPI.AddColumn(
		ctx, s.DB,
		&Table{Name: "compose_module_field"},
		&Column{Type: ColumnType{Type: ColumnTypeJson}, DefaultValue: "'{}'", Name: "config"},
	)
}

func fix202209_dropObsoleteComposeModuleFields(ctx context.Context, s *Store) (err error) {
	s.log(ctx).Info("extending compose_module_field table with config column")
	return s.SchemaAPI.DropColumn(
		ctx, s.DB,
		&Table{Name: "compose_module_field"},
		"is_private",
		"is_visible",
	)
}

func fix202209_composeRecordRevisions(ctx context.Context, s *Store) (err error) {
	s.log(ctx).Info("extending compose_record table with revision column")
	return s.SchemaAPI.AddColumn(
		ctx, s.DB,
		&Table{Name: "compose_record"},
		&Column{Type: ColumnType{Type: ColumnTypeInteger}, DefaultValue: "1", Name: "revision"},
	)

	//return // will probably not need this...
	//var (
	//	log = s.log(ctx)
	//)
	//
	//const (
	//	envKeySkip = "UPGRADE_COMPOSE_RECORD_CHANGES_PREFILL_SKIP"
	//
	//	tblChanges = "compose_record_changes"
	//	tblRecords = "compose_record"
	//)
	//
	//if _, set := os.LookupEnv(envKeySkip); set {
	//	return
	//}
	//
	//var (
	//	tblRecordColumns = []any{
	//		// reusing record ID for ID of the (1st record)
	//		// entry in the changes table
	//		"id",
	//		"id",
	//		"rel_namespace", "module_id",
	//		"values",
	//		"owned_by",
	//		"created_at", "created_by",
	//		"updated_at", "updated_by",
	//		"deleted_at", "deleted_by",
	//	}
	//
	//	tblChangesColumns = []any{
	//		"id",
	//		"rel_record",
	//		"rel_namespace", "rel_module",
	//		"values",
	//		"owned_by",
	//		"created_at", "created_by",
	//		"updated_at", "updated_by",
	//		"deleted_at", "deleted_by",
	//	}
	//
	//	d     = s.Dialect
	//	check = d.Select(exp.NewLiteralExpression("1")).From(tblChanges).Limit(1)
	//	count int
	//
	//	copyAll = d.
	//		Insert(tblChanges).
	//		Cols(tblChangesColumns...).
	//		FromQuery(d.From(tblRecords).Select(tblRecordColumns...))
	//)
	//
	//if err = s.QueryOne(ctx, check, &count); err != nil {
	//	// exit on error or when changes table is not empty
	//	if !errors.IsNotFound(err) {
	//		return fmt.Errorf("could not check if %q is empty: %w", tblChanges, err)
	//	}
	//
	//} else if count > 0 {
	//	return
	//}
	//
	//log.Warn(fmt.Sprintf("Empty %s table detected. "+
	//	"Prefilling record changes table with current records, "+
	//	"this might take a while, depending amount of records you have and how fast your database is. "+
	//	"If you need to disable this, stop the server, set %s=1 "+
	//	"and restart the server.",
	//	tblChanges,
	//	envKeySkip,
	//))
	//
	//spew.Dump(copyAll.ToSQL())
	//
	//// make a copy of records to changes table
	//return s.Exec(ctx, copyAll)
}
