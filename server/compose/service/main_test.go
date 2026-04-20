package service

import (
	"github.com/crusttech/human/server/pkg/cli"
	"github.com/crusttech/human/server/pkg/id"
)

func init() {
	id.Init(cli.Context())
}
