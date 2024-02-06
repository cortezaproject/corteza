package wfexec

import (
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

func init() {
	id.Init(cli.Context())
}
