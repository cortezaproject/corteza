package service

import (
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/revisions"
)

func diffResourceState(old, updated any) []*revisions.Change {
	return actionlog.DiffResourceState(old, updated)
}
