package event

import (
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &messageBase{
			channel: &types.Channel{Name: "ChanChan"},
			message: &types.Message{Message: "foo bar", Type: types.MessageTypeIlleism},
		}

		cMsg = eventbus.MustMakeConstraint("message", "like", "foo*")
		cTyp = eventbus.MustMakeConstraint("message.type", "eq", "illeism")
		cChn = eventbus.MustMakeConstraint("channel", "eq", "ChanChan")
	)

	a.True(messageMatch(res.message, cMsg))

	a.True(res.Match(cMsg))
	a.True(res.Match(cTyp))
	a.True(res.Match(cChn))
}
