package event

import (
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &pageBase{
			page:      &types.Page{Handle: "ph1"},
			namespace: &types.Namespace{Slug: "slg1"},
		}

		cPge = eventbus.MustMakeConstraint("page", "eq", "ph1")
		cNms = eventbus.MustMakeConstraint("namespace", "eq", "slg1")
	)

	a.True(res.Match(cPge))
	a.True(res.Match(cNms))
}
