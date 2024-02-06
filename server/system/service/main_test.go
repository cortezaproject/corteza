package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/id"
)

func init() {
	id.Init(context.Background())
}
