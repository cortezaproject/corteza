package rdbms

import (
	"context"

	systemType "github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) ApplicationMetrics(ctx context.Context) (_ *systemType.ApplicationMetrics, err error) {
	var (
		aux = struct {
			Total   uint `db:"total"`
			Deleted uint `db:"deleted"`
			Valid   uint `db:"valid"`
		}{}

		query = applicationSelectQuery(s.config.Dialect).
			Select(timestampStatExpr("deleted")...)
	)

	if err = s.QueryOne(ctx, query, &aux); err != nil {
		return nil, err
	}

	return &systemType.ApplicationMetrics{
		Total:   aux.Total,
		Deleted: aux.Deleted,
		Valid:   aux.Valid,
	}, nil
}

func (s Store) ReorderApplications(ctx context.Context, order []uint64) (err error) {
	//var (
	//	apps   systemType.ApplicationSet
	//	appMap = map[uint64]bool{}
	//	weight = 1
	//
	//	f = systemType.ApplicationFilter{}
	//)

	//s.Query(ctx, applicationSelectQuery())
	//if err = s.applicationCollection(ctx).Find().All(apps); err != nil {
	//	return
	//}
	//
	//if apps, _, err = s.SearchApplications(ctx, f); err != nil {
	//	return
	//}
	//
	//for _, app := range apps {
	//	appMap[app.ID] = true
	//}
	//
	//// honor parameter first
	//for _, id := range order {
	//	if appMap[id] {
	//		appMap[id] = false
	//
	//		app := apps.FindByID(id)
	//		if app == nil {
	//			continue
	//		}
	//
	//		app.Weight = weight
	//		weight++
	//	}
	//}
	//
	//for id, update := range appMap {
	//	if !update {
	//		continue
	//	}
	//
	//	app := apps.FindByID(id)
	//	if app == nil {
	//		continue
	//	}
	//
	//	app.Weight = weight
	//	weight++
	//}
	//
	//for _, app := range apps {
	//	err = s.applicationCollection(ctx).
	//		Find(db.Cond{"id": app.ID}).
	//		Update(app)
	//
	//	if err != nil {
	//		return
	//	}
	//}
	//
	return
}
