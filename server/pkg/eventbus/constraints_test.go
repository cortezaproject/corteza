package eventbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstraintMaking(t *testing.T) {
	var (
		a   = assert.New(t)
		c   ConstraintMatcher
		err error
	)

	c, err = ConstraintMaker("", "foo")
	a.Nil(c)
	a.Error(err)

	c, err = ConstraintMaker("", "eq")
	a.NoError(err)
	a.NotNil(c)
	a.Implements((*ConstraintMatcher)(nil), c)

	a.NotNil(MustMakeConstraint("foo", "=", "bar"))
}

func TestMustBeEqual(t *testing.T) {
	var (
		a      = assert.New(t)
		c, err = MustBeEqual("", "foo")
	)

	a.NoError(err)
	a.NotNil(c)
	a.True(c.Match("foo"))
	a.False(c.Match("bar"))
}

func TestMustNotBeEqual(t *testing.T) {
	var (
		a      = assert.New(t)
		c, err = MustNotBeEqual("", "foo")
	)

	a.NoError(err)
	a.NotNil(c)
	a.False(c.Match("foo"))
	a.True(c.Match("bar"))
}

func TestMustBeLike(t *testing.T) {
	var (
		a      = assert.New(t)
		c, err = MustBeLike("", "foo*")
	)

	a.NoError(err)
	a.NotNil(c)
	a.True(c.Match("fooBAZ"))
	a.False(c.Match("barBAZ"))
}

func TestMustNotBeLike(t *testing.T) {
	var (
		a      = assert.New(t)
		c, err = MustNotBeLike("", "foo*")
	)

	a.NoError(err)
	a.NotNil(c)
	a.False(c.Match("fooBAZ"))
	a.True(c.Match("barBAZ"))
}

func TestMustMatch(t *testing.T) {
	var (
		a      = assert.New(t)
		c, err = MustMatch("", "foo.+")
	)

	a.NoError(err)
	a.NotNil(c)
	a.True(c.Match("fooBAZ"))
	a.False(c.Match("barBAZ"))
}

func TestMustNotMatch(t *testing.T) {
	var (
		a      = assert.New(t)
		c, err = MustNotMatch("", "foo.*")
	)

	a.NoError(err)
	a.NotNil(c)
	a.False(c.Match("fooBAZ"))
	a.True(c.Match("barBAZ"))
}
