package outgoing

type (
	PayloadType interface {
		valid() bool
	}

	Error struct {
		Message string `json:"m"`
	}

	Message struct {
		ID        string `json:"id"`
		ChannelID string `json:"cid""`
		Message   string `json:"m"`
		Type      string `json:"t"`
		ReplyTo   string `json:"rid"`
		UserID    string `json:"uid"`
	}
)

func (*Error) valid() bool   { return true }
func (*Message) valid() bool { return true }
