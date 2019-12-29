package corredor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManualScriptFilterResourcePrefixing(t *testing.T) {
	f := &ManualScriptFilter{
		ResourceTypes: []string{"system", "system:one", "two"},
	}

	f.PrefixResource("system")

	assert.New(t).Equal(f.ResourceTypes, []string{"system", "system:one", "system:two"})
}
