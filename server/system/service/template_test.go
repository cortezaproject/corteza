package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFirstTemplateID(t *testing.T) {
	var req = require.New(t)

	// per-render override wins over the one set on the template
	req.Equal(uint64(1), firstTemplateID(1, 2))

	// no override, template's own header/footer is used
	req.Equal(uint64(2), firstTemplateID(0, 2))

	// neither set
	req.Zero(firstTemplateID(0, 0))
}
