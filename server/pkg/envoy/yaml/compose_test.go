package yaml

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestCompose_UnmarshalYAML(t *testing.T) {
	var (
		parse = func(c *compose, src string) (*compose, error) {
			return c, yaml.Unmarshal([]byte(src), c)
		}
	)

	t.Run("empty", func(t *testing.T) {
		req := require.New(t)
		c, err := parse(&compose{}, ``)
		req.NoError(err)
		req.NotNil(c)
		req.Empty(c.Namespaces)
	})

	t.Run("sequence of namespaces", func(t *testing.T) {
		req := require.New(t)
		c, err := parse(&compose{}, `namespaces: [ { } ]`)
		req.NoError(err)
		req.NotNil(c)
		req.Len(c.Namespaces, 1)
	})

	t.Run("map of namespaces", func(t *testing.T) {
		req := require.New(t)
		c, err := parse(&compose{}, `namespaces: { test: {} }`)
		req.NoError(err)
		req.NotNil(c)
		req.Len(c.Namespaces, 1)

	})

	t.Run("malformed namespace ref", func(t *testing.T) {
		req := require.New(t)
		_, err := parse(&compose{}, `namespace: { }`)
		req.Error(err)
	})

	t.Run("namespace ref", func(t *testing.T) {
		req := require.New(t)
		c, err := parse(&compose{}, `namespace: foo`)
		req.NoError(err)
		req.NotNil(c)
		req.Empty(c.Namespaces, "namespace ref should not result in namespace definition")
	})
}
