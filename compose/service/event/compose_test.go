package event

import (
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeOnIntervalMatching(t *testing.T) {
	var (
		res  = &composeOnInterval{}
		cInt = eventbus.MustMakeConstraint("", "", "* * * * *")
	)

	// Just make sure it runs
	// not bothering with precise test setup (scheduler's tests cover all that
	res.Match(cInt)
}

func TestComposeOnTimestampMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &composeOnInterval{}

		cTStamp = eventbus.MustMakeConstraint("", "", "2000-01-01T00:00:00Z") // Y2k!
	)

	a.False(res.Match(cTStamp), "Year 2000?! How did we get here? Anyhow, happy New year!")
}

func TestComposeMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &composeBase{}
	)

	a.False(res.Match(eventbus.MustMakeConstraint("foo", "", "bar")))
}
