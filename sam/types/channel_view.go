package types

type (
	ChannelView struct {
		ChannelID     uint64 `db:"rel_channel"`
		UserID        uint64 `db:"rel_user"`
		LastMessageID uint64 `db:"rel_last_message_id"`

		NewMessagesCount uint32 `db:"new_messages_count"`
	}

	ChannelViewFilter struct {
		UserID uint64
	}

	ChannelViewSet []*ChannelView
)

func (mm ChannelViewSet) Walk(w func(*ChannelView) error) (err error) {
	for i := range mm {
		if err = w(mm[i]); err != nil {
			return
		}
	}

	return
}

func (uu ChannelViewSet) FindByChannelId(channelID uint64) *ChannelView {
	for i := range uu {
		if uu[i].ChannelID == channelID {
			return uu[i]
		}
	}

	return nil
}
