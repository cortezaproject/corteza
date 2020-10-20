package yaml

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestComposeModule_UnmarshalYAML(t *testing.T) {
	var (
		req = require.New(t)

		w   *ComposeModule
		err error

		parse = func(w *ComposeModule, src string) (*ComposeModule, error) {
			return w, yaml.Unmarshal([]byte(src), w)
		}
	)

	w, err = parse(&ComposeModule{}, ``)
	req.NoError(err)
	req.NotNil(w)
	req.Nil(w.res)

	w, err = parse(&ComposeModule{}, `{ name: Test }`)
	req.NoError(err)
	req.NotNil(w)
	req.NotNil(w.res)
	req.NotEmpty(w.res.Name)
}
