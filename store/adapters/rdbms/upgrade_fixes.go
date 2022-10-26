package rdbms

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/model"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	labelsType "github.com/cortezaproject/corteza-server/pkg/label/types"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"strings"
	"time"
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
	fixes = []func(context.Context, *Store) error{
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
	}
)

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
		log = s.log(ctx)
	)
	_, err = s.DataDefiner.TableLookup(ctx, model.Record.Ident)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	const (
		recordSliceSize = 1000

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

	modules, _, err = s.SearchComposeModules(ctx, types.ModuleFilter{Deleted: filter.StateInclusive})
	if err != nil {
		return
	}

	log.Info(
		"preparing to migrate record values",
		zap.Int("modules", len(modules)),
		zap.Int("records", totalRecords),
	)

	// iterate through modules
	for _, mod := range modules {
		fields, _, err = s.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
		if err != nil {
			return
		}

		perModLog := log.With(
			zap.String("handle", mod.Handle),
			zap.Uint64("id", mod.ID),
		)

		err = func() (err error) {
			sliceLastRecordID = 0

			for {
				bmStart := time.Now()
				values = make(map[uint64]map[string][]any, recordSliceSize)
				recordIDs = make([]string, 0, recordSliceSize)

				err = func() (err error) {
					query = fmt.Sprintf(recordsPerModule, mod.NamespaceID, mod.ID, sliceLastRecordID, recordSliceSize)
					//println(query)
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
					//println(query)
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
