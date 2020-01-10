package event

import (
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemOnIntervalMatching(t *testing.T) {
	var (
		res  = &systemOnInterval{}
		cInt = eventbus.MustMakeConstraint("", "", "* * * * *")
	)

	// Just make sure it runs
	// not bothering with precise test setup (scheduler's tests cover all that
	res.Match(cInt)
}

func TestSystemOnTimestampMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &systemOnInterval{}

		cTStamp = eventbus.MustMakeConstraint("", "", "2000-01-01T00:00:00Z") // Y2k!
	)

	a.False(res.Match(cTStamp), "Year 2000?! How did we get here? Anyhow, happy New year!")
}

func TestSystemMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &systemBase{}
	)

	a.False(res.Match(eventbus.MustMakeConstraint("foo", "", "bar")))
}
