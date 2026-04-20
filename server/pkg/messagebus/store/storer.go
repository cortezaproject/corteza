package store

import (
	"github.com/crusttech/human/server/pkg/messagebus/types"
)

type (
	Storer interface {
		SetStore(types.QueueStorer)
		GetStore() types.QueueStorer
	}
)
