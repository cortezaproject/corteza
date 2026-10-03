package yaml

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// Language tags from imported YAML are capped like tags from requests, so
// an oversized tag cannot be used to exhaust the CPU through the parser
func TestDecodeLocale_oversizedTag(t *testing.T) {
	req := require.New(t)

	oversized := "en-" + strings.Repeat("x", 500)
	doc := "locale:\n  " + oversized + ":\n    key: value\n  en:\n    key: value\n"

	var n yaml.Node
	req.NoError(yaml.Unmarshal([]byte(doc), &n))

	rr, err := decodeLocale(n.Content[0])
	req.NoError(err)
	req.Len(rr, 2)

	req.Equal(language.Und, rr[0].locales[0].Lang.Tag, "oversized tag is replaced by und")
	req.Equal(language.English, rr[1].locales[0].Lang.Tag)
}
