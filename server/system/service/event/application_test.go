package event

import (
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/system/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplicationMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &applicationBase{
			application: &types.Application{Name: "someApp"},
		}

		cApp = eventbus.MustMakeConstraint("application.name", "eq", "someApp")
	)

	a.True(res.Match(cApp))
}
