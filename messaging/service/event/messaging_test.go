package event

import (
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessagingOnIntervalMatching(t *testing.T) {
	var (
		res  = &messagingOnInterval{}
		cInt = eventbus.MustMakeConstraint("", "", "* * * * *")
	)

	// Just make sure it runs
	// not bothering with precise test setup (scheduler's tests cover all that
	res.Match(cInt)
}

func TestMessagingOnTimestampMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &messagingOnInterval{}

		cTStamp = eventbus.MustMakeConstraint("", "", "2000-01-01T00:00:00Z") // Y2k!
	)

	a.False(res.Match(cTStamp), "Year 2000?! How did we get here? Anyhow, happy New year!")
}

func TestMessagingMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &messagingBase{}
	)

	a.False(res.Match(eventbus.MustMakeConstraint("foo", "", "bar")))
}
