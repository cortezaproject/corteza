package types

type (
	Unread struct {
		ChannelID     uint64 `db:"rel_channel"`
		ReplyTo       uint64 `db:"rel_reply_to"`
		UserID        uint64 `db:"rel_user"`
		LastMessageID uint64 `db:"rel_last_message"`

		Count uint32 `db:"count"`
	}

	UnreadFilter struct {
		UserID uint64
	}

	UnreadSet []*Unread
)

func (mm UnreadSet) Walk(w func(*Unread) error) (err error) {
	for i := range mm {
		if err = w(mm[i]); err != nil {
			return
		}
	}

	return
}

func (uu UnreadSet) FindByChannelId(channelID uint64) *Unread {
	for i := range uu {
		if uu[i].ChannelID == channelID {
			return uu[i]
		}
	}

	return nil
}
