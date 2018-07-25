package outgoing

type (
	Error struct {
		Message string `json:"m"`
	}

	Message struct {
		Id        string `json:"id"`
		ChannelId string `json:"cid""`
		Message   string `json:"m"`
		Type      string `json:"t"`
		ReplyTo   string `json:"rid"`
		UserId    string `json:"uid"`
	}
)
