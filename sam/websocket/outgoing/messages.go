package outgoing

type (
	Message struct {
		ID        string `json:"id"`
		ChannelID string `json:"cid"`
		Message   string `json:"m"`
		Type      string `json:"t"`
		ReplyTo   string `json:"rid"`
		UserID    string `json:"uid"`
	}

	Messages []*Message

	MessageUpdate struct {
		ID      string `json:"id"`
		Message string `json:"m"`
	}

	MessageDelete struct {
		ID string `json:"id"`
	}
)

func (*Message) valid() bool  { return true }
func (*Messages) valid() bool { return true }

func (*MessageUpdate) valid() bool { return true }
func (*MessageDelete) valid() bool { return true }
