package service

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

type noGlobalFlag struct{ allowAll }

func (noGlobalFlag) CanGlobalApplicationFlag(context.Context) bool { return false }

func TestClaimHomeClearsItElsewhere(t *testing.T) {
	ctx := context.Background()
	svc, first := newTestApplication(t)

	first.Unify.Home = true
	require.NoError(t, store.UpdateApplication(ctx, svc.store, first))

	second := &types.Application{ID: nextID(), Name: "second", CreatedAt: *now(), Unify: &types.ApplicationUnify{}}
	require.NoError(t, store.CreateApplication(ctx, svc.store, second))

	require.NoError(t, svc.claimHome(ctx, second.ID, false, true))

	first, err := store.LookupApplicationByID(ctx, svc.store, first.ID)
	require.NoError(t, err)
	require.False(t, first.Unify.Home)
}

func TestClaimHomeTakesTheGlobalFlagPermission(t *testing.T) {
	ctx := context.Background()
	svc, app := newTestApplication(t)
	svc.ac = noGlobalFlag{}

	require.ErrorIs(t, svc.claimHome(ctx, app.ID, false, true), ApplicationErrNotAllowedToManageFlagGlobal())
	require.ErrorIs(t, svc.claimHome(ctx, app.ID, true, false), ApplicationErrNotAllowedToManageFlagGlobal())
	require.NoError(t, svc.claimHome(ctx, app.ID, true, true))
}
