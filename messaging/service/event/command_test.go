package event

import (
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &commandBase{
			channel: &types.Channel{Name: "ChanChan"},
			command: &types.Command{Name: "fooCommand"},
		}

		cFoo = eventbus.MustMakeConstraint("command", "eq", "fooCommand")
		cBar = eventbus.MustMakeConstraint("command", "eq", "barCommand")
		cChn = eventbus.MustMakeConstraint("channel", "eq", "ChanChan")
	)

	a.True(commandMatch(res.command, cFoo))
	a.False(commandMatch(res.command, cBar))

	a.True(res.Match(cFoo))
	a.False(res.Match(cBar))
	a.True(res.Match(cChn))
}
