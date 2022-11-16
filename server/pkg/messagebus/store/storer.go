package store

import (
	"github.com/cortezaproject/corteza/server/pkg/messagebus/types"
)

type (
	Storer interface {
		SetStore(types.QueueStorer)
		GetStore() types.QueueStorer
	}
)
