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

	ChannelJoin struct {
		// ID of the channel user is joining
		ID string `json:"id"`

		// ID of the user that is joining
		UserID string `json:"uid"`
	}

	ChannelPart struct {
		// Channel to part (nil) for ALL channels
		ID *string `json:"id"`

		// Who is parting
		UserID string `json:"uid"`
	}

	Channel struct {
		// Channel to part (nil) for ALL channels
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	Channels []*Channel
)

func (*Message) valid() bool  { return true }
func (*Messages) valid() bool { return true }

func (*MessageUpdate) valid() bool { return true }
func (*MessageDelete) valid() bool { return true }

func (*ChannelJoin) valid() bool { return true }
func (*ChannelPart) valid() bool { return true }
func (*Channel) valid() bool     { return true }
func (*Channels) valid() bool    { return true }
