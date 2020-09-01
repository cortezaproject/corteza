package types

import (
	"time"
)

type (
	Mention struct {
		ID            uint64
		MessageID     uint64
		ChannelID     uint64
		UserID        uint64
		MentionedByID uint64
		CreatedAt     time.Time
	}

	MentionFilter struct {
		MessageID []uint64
		//// All mentions by this user
		//MentionedByID []uint64
		//
		//// All mentions of this user
		//UserID []uint64
		//
		//// How many entries
		//Limit uint
	}
)
