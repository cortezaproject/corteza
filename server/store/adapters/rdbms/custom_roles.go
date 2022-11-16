package rdbms

import (
	"context"

	systemType "github.com/cortezaproject/corteza/server/system/types"
	"github.com/doug-martin/goqu/v9"
)

func (s Store) RoleMetrics(ctx context.Context) (m *systemType.RoleMetrics, err error) {
	var (
		aux = struct {
			Total    uint `db:"total"`
			Deleted  uint `db:"deleted"`
			Valid    uint `db:"valid"`
			Archived uint `db:"archived"`
		}{}

		query = roleSelectQuery(s.Dialect.GOQU()).
			Select(timestampStatExpr("deleted", "archived")...)
	)

	if err = s.QueryOne(ctx, query, &aux); err != nil {
		return nil, err
	}

	m = &systemType.RoleMetrics{
		Total:    aux.Total,
		Valid:    aux.Valid,
		Deleted:  aux.Deleted,
		Archived: aux.Archived,
	}

	// Fetch daily metrics for created, updated, deleted and suspended users
	err = s.multiDailyMetrics(
		ctx,
		roleTable,
		[]goqu.Expression{},
		[]string{
			"created_at",
			"updated_at",
			"deleted_at",
			"archived_at",
		},
		&m.DailyCreated,
		&m.DailyUpdated,
		&m.DailyDeleted,
		&m.DailyArchived,
	)

	if err != nil {
		return nil, err
	}

	return
}

func (s Store) TransferRoleMembers(ctx context.Context, src, dst uint64) (err error) {
	var (
		transfer = s.Dialect.GOQU().Update(roleMemberTable).
			Set(goqu.Record{"rel_role": dst}).
			Where(goqu.Ex{"rel_role": src})
	)

	return s.Exec(ctx, transfer)
}
