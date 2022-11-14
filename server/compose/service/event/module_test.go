package event

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuleMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &moduleBase{
			module:    &types.Module{Handle: "mh1"},
			namespace: &types.Namespace{Slug: "slg1"},
		}

		cMod = eventbus.MustMakeConstraint("module", "eq", "mh1")
		cNms = eventbus.MustMakeConstraint("namespace", "eq", "slg1")
	)

	a.True(res.Match(cMod))
	a.True(res.Match(cNms))
}
