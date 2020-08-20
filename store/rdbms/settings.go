package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) columns() []string {
	return []string{
		"name",
		"value",
		"rel_owner",
		"updated_at",
		"updated_by",
	}
}

func (s Store) convertSettingFilter(f types.SettingsFilter) (query squirrel.SelectBuilder, err error) {
	query = s.settingsSelectBuilder().Where(squirrel.Eq{"rel_owner": f.OwnedBy})

	if len(f.Prefix) > 0 {
		query = query.Where("name LIKE ?", f.Prefix+"%")
	}

	return
}

//func (s Store) BulkSet(ctx context.Context, vv ValueSet) error {
//	// Save all inside a db transaction
//	return r.db().Transaction(func() (err error) {
//		return vv.Walk(func(v *SettingValue) error {
//			return r.Set(v)
//		})
//	})
//}
//
//func (s Store) UpdateSettings(ctx context.Context, value *SettingValue) error {
//	return s.DB().Replace(s.dbTable, value)
//}
//
//func (s Store) DeleteSettings(name string, ownedBy uint64) error {
//	_, err := s.db().Exec(
//		fmt.Sprintf("DELETE FROM %s WHERE name = ? AND rel_owner = ? ", s.dbTable),
//		name,
//		ownedBy,
//	)
//	return err
//}
//
//func (s Store) SettingsLookupByName(ctx context.Context, name string, ownedBy uint64) (*settings.SettingValue, error) {
//	var (
//		v   = &settings.SettingValue{}
//		cnd = squirrel.Eq{"rel_owner": ownedBy, "name": name}
//		err = s.Lookup(ctx, v, s.QuerySettings(), cnd)
//	)
//
//	if err != nil {
//		return
//	}
//
//	return v, err
//}
//
//// Returns squirrel.selectBuilder
//func (s Store) QuerySettings() squirrel.selectBuilder {
//	return s.selectBuilder(s.SettingsTable()+" AS stngs", s.SettingsColumns()...)
//}
//
//// Name of the db table
//func (Store) SettingsTable() string {
//	return "sys_settings"
//}
//
//func (Store) SettingsColumns() []string {
//	return []string{
//		"stngs.name",
//		"stngs.value",
//		"stngs.rel_owner",
//		"stngs.created_by",
//		"stngs.updated_at",
//	}
//}
