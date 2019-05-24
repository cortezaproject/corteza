package compose

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/internal/test"
)

func TestConfigure(t *testing.T) {
	var config = Configure()
	test.Assert(t, config != nil, "Configure valid")
	test.Assert(t, func() bool { config.Init(); return true }(), "Initialization ok")
	test.Assert(t, config.MakeCLI(context.Background()) != nil, "CLI created")
}
