package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPayload_Only(t *testing.T) {
	a := assert.New(t)

	p := Payload{"foo": 1, "bar": 1}.Only("foo")
	a.Contains(p, "foo")
	a.NotContains(p, "bar")
}

func TestPayload_Skip(t *testing.T) {
	a := assert.New(t)

	p := Payload{"foo": 1, "bar": 1}.Skip("bar")
	a.Contains(p, "foo")
	a.NotContains(p, "bar")
}
