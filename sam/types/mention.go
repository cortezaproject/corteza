package types

import (
	"time"
)

type (
	Mention struct {
		ID            uint64    `db:"id"`
		MessageID     uint64    `db:"rel_message"`
		ChannelID     uint64    `db:"rel_channel"`
		UserID        uint64    `db:"rel_user"`
		MentionedByID uint64    `db:"rel_mentioned_by"`
		CreatedAt     time.Time `db:"created_at"`
	}

	MentionSet []*Mention

	MentionFilter struct {
		// All mentions by this user
		MentionedByID uint64

		// All mentions of this user
		UserID uint64

		// How many entries
		Limit uint
	}
)

func (mm MentionSet) Walk(w func(*Mention) error) (err error) {
	for i := range mm {
		if err = w(mm[i]); err != nil {
			return
		}
	}

	return
}

func (mm MentionSet) FindByID(ID uint64) (out *Mention) {
	out = &Mention{}

	for i := range mm {
		if mm[i].ID == ID {
			return
		}
	}

	return nil
}

func (mm MentionSet) FindByUserID(ID uint64) (out MentionSet) {
	out = MentionSet{}

	for i := range mm {
		if mm[i].UserID == ID {
			out = append(out, mm[i])
		}
	}

	return
}

func (mm MentionSet) FindByMessageID(ID uint64) (out MentionSet) {
	out = MentionSet{}

	for i := range mm {
		if mm[i].MessageID == ID {
			out = append(out, mm[i])
		}
	}

	return
}

func (mm MentionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(mm))

	for i := range mm {
		IDs[i] = mm[i].ID
	}

	return
}

func (mm MentionSet) UserIDs() (IDs []uint64) {
	IDs = make([]uint64, len(mm))

	for i := range mm {
		IDs[i] = mm[i].UserID
	}

	return
}

func (mm MentionSet) Diff(in MentionSet) (add, upd, del MentionSet) {
	add, upd, del = MentionSet{}, MentionSet{}, MentionSet{}

	for _, m := range in {
		if m.ID == 0 {
			// Mark for adding all new
			add = append(add, m)
		}
	}

	for _, m := range mm {
		if m.ID == 0 {
			// Ignore all unsaved
			continue
		}

		if in.FindByID(m.ID) == nil {
			// Mark for removal all that are not added
			del = append(del, m)
		} else {
			// Mark for update all that are still there
			upd = append(upd, m)
		}
	}

	return
}
