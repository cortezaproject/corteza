package corredor

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type (
	mockEvent struct {
		rType string
		eType string
		match func(matcher eventbus.ConstraintMatcher) bool
	}

	mockUserSvc struct {
		user *types.User
		err  error
	}
	mockRoleSvc struct {
		role *types.Role
		err  error
	}
)

func (e mockEvent) ResourceType() string {
	return e.rType
}

func (e mockEvent) EventType() string {
	return e.eType
}

func (e mockEvent) Encode() (map[string][]byte, error) {
	return nil, nil
}

func (e mockEvent) Decode(map[string][]byte) error {
	return nil
}

func (e mockEvent) Match(matcher eventbus.ConstraintMatcher) bool {
	if e.match == nil {
		return true
	}

	return e.match(matcher)
}

func (u *mockUserSvc) FindByAny(context.Context, interface{}) (*types.User, error) {
	return u.user, u.err
}

func (u *mockRoleSvc) FindByAny(context.Context, interface{}) (*types.Role, error) {
	return u.role, u.err
}

func TestFindOnManual(t *testing.T) {
	var (
		ctx = context.Background()

		svc = &service{
			denyExec: make(map[string]map[uint64]bool),
			sScripts: ScriptSet{
				&Script{
					Name: "s1",
					Triggers: []*Trigger{
						&Trigger{
							EventTypes:    []string{"ev"},
							ResourceTypes: []string{"res"},
						},
					},
				},
				&Script{
					Name: "s2",
					Triggers: []*Trigger{
						&Trigger{
							EventTypes:    []string{"foo"},
							ResourceTypes: []string{"bar"},
						},
					},
				},
				//&Script{
				//	Triggers: []*Trigger{
				//		&Trigger{
				//			EventTypes:    []string{"not-a-match"},
				//			ResourceTypes: []string{"not-a-match"},
				//		},
				//	},
				//},
			},
			cScripts: ScriptSet{
				&Script{
					Name: "s3",
					Triggers: []*Trigger{
						&Trigger{
							EventTypes:    []string{"ev"},
							ResourceTypes: []string{"res"},
						},
					},
				},
				&Script{
					Name: "s4",
					Triggers: []*Trigger{
						&Trigger{
							EventTypes:    []string{"foo"},
							ResourceTypes: []string{"bar"},
						},
					},
				},
			},
		}
		filter = Filter{
			ResourceTypes:        []string{"res"},
			EventTypes:           []string{"ev"},
			ExcludeServerScripts: false,
			ExcludeClientScripts: false,
		}

		o, _, err = svc.Find(ctx, filter)

		a = assert.New(t)
	)

	a.NoError(err)
	a.Len(o, 2)
}

func TestGlobals(t *testing.T) {
	var (
		a = assert.New(t)
	)

	gCorredor = NewService(zap.NewNop(), options.CorredorOpt{})
	a.Equal(gCorredor, Service())
	a.NoError(Setup(zap.NewNop(), options.CorredorOpt{}))
	a.NotNil(gCorredor)

	gCorredor = nil
	a.NoError(Setup(zap.NewNop(), options.CorredorOpt{}))
	a.Equal(gCorredor, Service())
	a.NotNil(gCorredor)
}

func TestServiceBasics(t *testing.T) {
	var (
		svc = &service{}
	)

	svc.SetEventRegistry(nil)
	svc.SetAuthTokenMaker(nil)
	svc.SetUserFinder(nil)
	svc.SetRoleFinder(nil)
}

func TestService_canExec(t *testing.T) {
	var (
		a   = assert.New(t)
		svc = &service{
			users:    &mockUserSvc{user: &types.User{ID: 42, Email: "dummy@mo.ck", Handle: "dummy"}},
			roles:    &mockRoleSvc{role: &types.Role{ID: 84, Handle: "role", Name: "ROLE"}},
			denyExec: make(map[string]map[uint64]bool),
		}

		ctx = auth.SetIdentityToContext(context.Background(), auth.Authenticated(42, 84))

		script1 = &ServerScript{
			Name: "s1",
			Triggers: []*Trigger{&Trigger{
				EventTypes:    []string{onManualEventType},
				ResourceTypes: []string{"res"},
			}},
			Security: &Security{
				RunAs: "foo",
				Allow: []string{"role"},
			},
		}

		script2 = &ServerScript{
			Name: "s2",
			Triggers: []*Trigger{&Trigger{
				EventTypes:    []string{onManualEventType},
				ResourceTypes: []string{"res"},
			}},
			Security: &Security{
				RunAs: "foo",
				Deny:  []string{"role"},
			},
		}

		script3 = &ServerScript{
			Name: "s3", // permissions will not be added, not a manual script
			Security: &Security{
				RunAs: "foo",
				Deny:  []string{"role"},
			},
		}

		script4 = &ServerScript{
			Name: "s3", // should not be even added, duplicated name (hence length=3)
			Security: &Security{
				RunAs: "foo",
				Deny:  []string{"role"},
			},
		}
	)

	if testing.Verbose() {
		svc.log = logger.MakeDebugLogger()
	} else {
		svc.log = zap.NewNop()
	}

	svc.registerServerScripts(ctx, script1, script2, script3, script4)

	a.Len(svc.sScripts, 3)
	a.Len(svc.denyExec, 2)

	a.True(svc.canExec(ctx, script1.Name))
	a.False(svc.canExec(ctx, script2.Name))
}

func TestService_Exec(t *testing.T) {
	var (
		a   = assert.New(t)
		svc = &service{
			explicit: map[string]map[string]bool{},
			opt:      options.CorredorOpt{Enabled: true},
		}
	)

	a.Error(svc.Exec(nil, "script", &mockEvent{eType: "not-onManual"}))
	a.Error(svc.Exec(nil, "script", &mockEvent{eType: onManualEventType}))
	svc.explicit["script"] = nil
	a.Error(svc.Exec(nil, "script", &mockEvent{eType: onManualEventType}))
}
