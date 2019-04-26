package types

type (
	Activity struct {
		ChannelID uint64 `json:"channelID,string"`
		MessageID uint64 `json:"messageID,string"`
		UserID    uint64 `json:"userID,string"`
		Kind      string `json:"kind,omitempty"`
	}
)
