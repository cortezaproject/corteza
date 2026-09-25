package service

import (
	"context"
	"testing"

	internalAuth "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

type selfFlagOnly struct {
	userAccessController
	allowed bool
}

func (ac selfFlagOnly) CanSelfApplicationFlag(context.Context) bool { return ac.allowed }

func TestKeepHomeApplication(t *testing.T) {
	const self, other uint64 = 1, 2
	ctx := internalAuth.SetIdentityToContext(context.Background(), internalAuth.Authenticated(self))

	pick := func(allowed bool, userID uint64) uint64 {
		svc := &user{ac: selfFlagOnly{allowed: allowed}}
		res := &types.User{ID: userID, Meta: &types.UserMeta{HomeApplicationID: 10}}
		upd := &types.User{ID: userID, Meta: &types.UserMeta{HomeApplicationID: 20}}
		svc.keepHomeApplication(ctx, upd, res)
		return upd.Meta.HomeApplicationID
	}

	require.Equal(t, uint64(20), pick(true, self), "own pick with the permission")
	require.Equal(t, uint64(10), pick(false, self), "own pick without it is kept as it was")
	require.Equal(t, uint64(20), pick(false, other), "updating another user is CanUpdateUser's call")
}
