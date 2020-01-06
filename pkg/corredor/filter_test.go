package corredor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManualScriptFilterResourcePrefixing(t *testing.T) {
	f := &Filter{
		ResourceTypes: []string{"system", "system:one", "two"},
	}

	f.PrefixResources("system")

	assert.New(t).Equal(f.ResourceTypes, []string{"system", "system:one", "system:two"})
}
