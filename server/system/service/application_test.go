package service

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type allowAll struct{}

func (allowAll) CanCreateApplication(context.Context) bool                             { return true }
func (allowAll) CanSearchApplications(context.Context) bool                            { return true }
func (allowAll) CanSelfApplicationFlag(context.Context) bool                           { return true }
func (allowAll) CanGlobalApplicationFlag(context.Context) bool                         { return true }
func (allowAll) CanReadApplication(context.Context, *types.Application) bool           { return true }
func (allowAll) CanUpdateApplication(context.Context, *types.Application) bool         { return true }
func (allowAll) CanDeleteApplication(context.Context, *types.Application) bool         { return true }
func (allowAll) CanManageSourceOnApplication(context.Context, *types.Application) bool { return true }

func newTestApplication(t *testing.T) (*application, *types.Application) {
	t.Helper()
	ctx := context.Background()

	s, err := sqlite.ConnectInMemory(ctx)
	require.NoError(t, err)
	require.NoError(t, store.Upgrade(ctx, zap.NewNop(), s))

	app := &types.Application{
		ID:        nextID(),
		Name:      "probe",
		CreatedAt: *now(),
		Unify:     &types.ApplicationUnify{Kind: ApplicationKindCustom},
	}
	require.NoError(t, store.CreateApplication(ctx, s, app))

	return &application{
		store:     s,
		ac:        allowAll{},
		actionlog: actionlog.NewService(s, zap.NewNop(), zap.NewNop(), actionlog.MakeDisabledPolicy()),
		services:  &applicationServices{eventbus: nil},
	}, app
}

// The caller resolves the declaration, because this package cannot reach
// compose; storing one that was not resolved would leave every viewer to
// search for what it names.
func TestSetSourceRefusesAnUnresolvedDeclaration(t *testing.T) {
	svc, app := newTestApplication(t)

	err := svc.onSetSource(context.Background(), &applicationActionProps{}, app, "<p>x</p>", &types.ApplicationSourceMeta{
		Namespace: "crm",
		Modules:   []string{"Lead"},
	})
	require.ErrorIs(t, err, ApplicationErrDeclarationNotResolved())

	err = svc.onSetSource(context.Background(), &applicationActionProps{}, app, "<p>x</p>", &types.ApplicationSourceMeta{
		Namespace:   "crm",
		Modules:     []string{"Lead"},
		NamespaceID: 42,
		ModuleIDs:   map[string]string{"Lead": "43"},
	})
	require.NoError(t, err)
}
