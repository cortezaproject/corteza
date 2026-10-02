package mcp_test

import (
	"fmt"
	"testing"

	cmpService "github.com/crusttech/human/server/compose/service"
	herrors "github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/stretchr/testify/assert"
)

// TestClassifyError pins how Human's own errors read on the wire. The
// generated action errors are all KindInternal, so the classifier has to read
// their "type" meta; a wrong mapping here sends a model the wrong next step.
func TestClassifyError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		code string
	}{
		{"generated not found", cmpService.NamespaceErrNotFound(), toolkit.CodeNotFound},
		{"wrapped by a handler", fmt.Errorf("namespace lookup failed: %w", cmpService.ModuleErrNotFound()), toolkit.CodeNotFound},
		{"generated not allowed", cmpService.ModuleErrNotAllowedToUpdate(), toolkit.CodeForbidden},
		{"generated invalid id", cmpService.PageErrInvalidID(), toolkit.CodeInvalid},
		{"generated stale data", cmpService.ModuleErrStaleData(), toolkit.CodeConflict},
		{"kind not found", herrors.New(herrors.KindNotFound, "gone"), toolkit.CodeNotFound},
		{"kind unauthorized", herrors.New(herrors.KindUnauthorized, "no"), toolkit.CodeForbidden},
		{"kind invalid data", herrors.New(herrors.KindInvalidData, "bad"), toolkit.CodeInvalid},
		{"unknown", fmt.Errorf("disk on fire"), ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, _ := mcp.ClassifyError(c.err)
			assert.Equal(t, c.code, code)
		})
	}
}
