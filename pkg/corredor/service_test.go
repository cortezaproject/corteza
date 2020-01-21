package corredor

import (
	"github.com/cortezaproject/corteza-server/pkg/app/options"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"go.uber.org/zap"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	mockEvent struct {
		rType string
		eType string
		match func(matcher eventbus.ConstraintMatcher) bool
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

func TestFindOnManual(t *testing.T) {
	var (
		svc = &service{
			sScripts: ScriptSet{
				&Script{
					Triggers: []*Trigger{
						&Trigger{
							EventTypes:    []string{"ev"},
							ResourceTypes: []string{"res"},
						},
					},
				},
				&Script{
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
					Triggers: []*Trigger{
						&Trigger{
							EventTypes:    []string{"ev"},
							ResourceTypes: []string{"res"},
						},
					},
				},
				&Script{
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

		o, _, err = svc.Find(filter)

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
}

func TestService_ExecOnManual(t *testing.T) {
	var (
		a   = assert.New(t)
		svc = &service{
			manual: map[string]map[string]string{},
		}
	)

	a.Error(svc.ExecOnManual(nil, "script", &mockEvent{eType: "not-onManual"}))
	a.Error(svc.ExecOnManual(nil, "script", &mockEvent{eType: onManualEventType}))
	svc.manual["script"] = nil
	a.Error(svc.ExecOnManual(nil, "script", &mockEvent{eType: onManualEventType}))
}
