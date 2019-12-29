package corredor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindOnManual(t *testing.T) {
	var (
		svc = &service{
			sScripts: ScriptSet{
				&Script{
					Triggers: []*Trigger{
						&Trigger{
							Events:    []string{"ev"},
							Resources: []string{"res"},
						},
					},
				},
				&Script{
					Triggers: []*Trigger{
						&Trigger{
							Events:    []string{"foo"},
							Resources: []string{"bar"},
						},
					},
				},
			},
			cScripts: ScriptSet{
				&Script{
					Triggers: []*Trigger{
						&Trigger{
							Events:    []string{"ev"},
							Resources: []string{"res"},
						},
					},
				},
				&Script{
					Triggers: []*Trigger{
						&Trigger{
							Events:    []string{"foo"},
							Resources: []string{"bar"},
						},
					},
				},
			},
		}
		filter = ManualScriptFilter{
			ResourceTypes:        []string{"res"},
			EventTypes:           []string{"ev"},
			ExcludeServerScripts: false,
			ExcludeClientScripts: false,
		}

		o, _, err = svc.FindOnManual(filter)

		a = assert.New(t)
	)

	a.NoError(err)
	a.Len(o, 2)
}
