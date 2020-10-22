package yaml

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestCompose_UnmarshalYAML(t *testing.T) {
	var (
		req = require.New(t)

		c   *compose
		err error

		parse = func(c *compose, src string) (*compose, error) {
			return c, yaml.Unmarshal([]byte(src), c)
		}
	)

	c, err = parse(&compose{}, ``)
	req.NoError(err)
	req.NotNil(c)
	req.Empty(c.namespaces)

	c, err = parse(&compose{}, `namespaces: [ { } ]`)
	req.NoError(err)
	req.NotNil(c)
	req.Len(c.namespaces, 1)

	c, err = parse(&compose{}, `namespaces: { test: {} }`)
	req.NoError(err)
	req.NotNil(c)
	req.Len(c.namespaces, 1)

	c, err = parse(&compose{}, `namespace: { }`)
	req.Error(err)

	c, err = parse(&compose{}, `namespace: foo`)
	req.NoError(err)
	req.NotNil(c)
	req.Empty(c.namespaces, "namespace ref should not result in namespace definition")
}
