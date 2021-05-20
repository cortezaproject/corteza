package rdbms

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
)

func (s Store) TransferRbacRules(ctx context.Context, src, dst uint64) (err error) {
	return s.execUpdateRbacRules(ctx, squirrel.Eq{"rel_role": src}, store.Payload{"rel_role": dst})
}
