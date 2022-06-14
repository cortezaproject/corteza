package dalutils

import (
	"context"
	"sort"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	sensitivityLevelCreator interface {
		CreateSensitivityLevel(levels ...dal.SensitivityLevel) (err error)
	}

	sensitivityLevelUpdater interface {
		UpdateSensitivityLevel(levels ...dal.SensitivityLevel) (err error)
	}

	sensitivityLevelDeleter interface {
		DeleteSensitivityLevel(levels ...dal.SensitivityLevel) (err error)
	}

	sensitivityLevelReloader interface {
		ReloadSensitivityLevels(levels ...dal.SensitivityLevel) (err error)
	}
)

func DalSensitivityLevelReload(ctx context.Context, s store.Storer, r sensitivityLevelReloader) (err error) {
	ll, _, err := store.SearchDalSensitivityLevels(ctx, s, types.DalSensitivityLevelFilter{Deleted: filter.StateExcluded})
	if err != nil {
		return
	}

	sort.Sort(ll)

	levels := make(dal.SensitivityLevelSet, 0, len(ll))
	for _, l := range ll {
		levels = append(levels, systemToPkgType(l))
	}

	return r.ReloadSensitivityLevels(levels...)
}

func DalSensitivityLevelCreate(c sensitivityLevelCreator, levels ...*types.DalSensitivityLevel) (err error) {
	return c.CreateSensitivityLevel(systemToPkgTypeSet(levels)...)
}

func DalSensitivityLevelUpdate(u sensitivityLevelUpdater, levels ...*types.DalSensitivityLevel) (err error) {
	return u.UpdateSensitivityLevel(systemToPkgTypeSet(levels)...)
}

func DalSensitivityLevelDelete(d sensitivityLevelDeleter, levels ...*types.DalSensitivityLevel) (err error) {
	return d.DeleteSensitivityLevel(systemToPkgTypeSet(levels)...)
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utils

func systemToPkgType(level *types.DalSensitivityLevel) dal.SensitivityLevel {
	return dal.SensitivityLevel{
		ID:     level.ID,
		Handle: level.Handle,
		Level:  level.Level,
	}
}

func systemToPkgTypeSet(levels types.DalSensitivityLevelSet) dal.SensitivityLevelSet {
	out := make(dal.SensitivityLevelSet, 0, len(levels))
	for _, l := range levels {
		out = append(out, systemToPkgType(l))
	}
	return out
}
