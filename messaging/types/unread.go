package types

type (
	Unread struct {
		ChannelID     uint64
		ReplyTo       uint64
		UserID        uint64
		LastMessageID uint64

		Count       uint32
		ThreadCount uint32
		ThreadTotal uint32
	}
)

func (uu UnreadSet) Merge(in UnreadSet) UnreadSet {
	var (
		out  = append(UnreadSet{}, uu...)
		olen = len(out)
	)

inSet:
	for _, i := range in {
		for o := 0; o < olen; o++ {
			if out[o].UserID == i.UserID && out[o].ChannelID == i.ChannelID && out[o].ReplyTo == i.ReplyTo {
				if i.Count > 0 {
					out[o].Count = i.Count
				}
				if i.ThreadCount > 0 {
					out[o].ThreadCount = i.ThreadCount
				}
				if i.ThreadTotal > 0 {
					out[o].ThreadTotal = i.ThreadTotal
				}

				continue inSet
			}
		}

		out = append(out, i)
	}

	return out
}
