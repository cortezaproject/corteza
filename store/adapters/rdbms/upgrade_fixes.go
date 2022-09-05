package rdbms

//func fix202209_extendComposeModuleForPrivacyAndDAL(ctx context.Context, s *Store) (err error) {
//	s.log(ctx).Info("extending compose_module table with config column")
//	return s.SchemaAPI.AddColumn(
//		ctx, s.DB,
//		"compose_module",
//		&Column{Type: ColumnType{Type: ColumnTypeJson}, DefaultValue: "'{}'", Name: "config"},
//	)
//}
//
//func fix202209_extendComposeModuleFieldsForPrivacyAndDAL(ctx context.Context, s *Store) (err error) {
//	s.log(ctx).Info("extending compose_module_field table with config column")
//	return s.SchemaAPI.AddColumn(
//		ctx, s.DB,
//		"compose_module_field",
//		&Column{Type: ColumnType{Type: ColumnTypeJson}, DefaultValue: "'{}'", Name: "config"},
//	)
//}
//
//func fix202209_dropObsoleteComposeModuleFields(ctx context.Context, s *Store) (err error) {
//	s.log(ctx).Info("extending compose_module_field table with config column")
//	return s.SchemaAPI.DropColumn(
//		ctx, s.DB,
//		"compose_module_field",
//		"is_private",
//		"is_visible",
//	)
//}
//
//func fix202209_composeRecordRevisions(ctx context.Context, s *Store) (err error) {
//	s.log(ctx).Info("extending compose_record table with revision column")
//	return s.SchemaAPI.AddColumn(
//		ctx, s.DB,
//		"compose_record",
//		&Column{Type: ColumnType{Type: ColumnTypeInteger}, DefaultValue: "1", Name: "revision"},
//	)

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
//}
//
//func fix202209_extendDalConnectionsForMeta(ctx context.Context, s *Store) (err error) {
//	s.log(ctx).Info("extending dal_connections table with meta column")
//	return s.SchemaAPI.AddColumn(
//		ctx, s.DB,
//		&Table{Name: "dal_connections"},
//		&Column{Type: ColumnType{Type: ColumnTypeJson}, DefaultValue: "'{}'", Name: "meta"},
//	)
//}
//
//func fix202209_renameModuleColOnComposeRecords(ctx context.Context, s *Store) (err error) {
//	s.log(ctx).Info("rename module_id column on compose_record table")
//	return s.SchemaAPI.RenameColumn(
//		ctx, s.DB, &Table{Name: "compose_record"}, "module_id", "rel_module",
//	)
//}
//
//func fix202209_addMetaOnComposeRecords(ctx context.Context, s *Store) (err error) {
//	var (
//		log = s.log(ctx)
//
//		groupedMeta = make(map[uint64]map[string]any)
//		packed      []byte
//	)
//
//	log.Info("add meta column on compose_record table")
//	err = s.SchemaAPI.AddColumn(
//		ctx, s.DB,
//		&Table{Name: "compose_record"},
//		&Column{Type: ColumnType{Type: ColumnTypeJson}, DefaultValue: "'{}'", Name: "meta"},
//	)
//
//	if err != nil {
//		return
//	}
//
//	return s.Tx(ctx, func(ctx context.Context, s store.Storer) (err error) {
//		log.Info("collecting record labels")
//		ll, _, err := store.SearchLabels(ctx, s, labelsType.LabelFilter{Kind: "compose:record"})
//		if err != nil {
//			return
//		}
//
//		log.Info("grouping labels", zap.Int("count", len(ll)))
//		for _, l := range ll {
//			if _, has := groupedMeta[l.ResourceID]; !has {
//				groupedMeta[l.ResourceID] = make(map[string]any)
//			}
//
//			groupedMeta[l.ResourceID][l.Name] = l.Value
//			if err = store.DeleteLabel(ctx, s, l); err != nil {
//				return
//			}
//		}
//
//		log.Info("updating records with meta", zap.Int("count", len(ll)))
//		for recordID, labels := range groupedMeta {
//			packed, err = json.Marshal(labels)
//			_, err = s.(*Store).DB.ExecContext(ctx, "UPDATE compose_record SET meta = $1 WHERE id = $2", packed, recordID)
//			if err != nil {
//				return
//			}
//		}
//
//		return
//
//	})
//
//}
