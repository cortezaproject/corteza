package outgoing

type (
	Unread struct {
		// Channel to part (nil) for ALL channels
		LastMessageID uint64 `json:"lastMessageID,string,omitempty"`
		Count         uint32 `json:"count"`
	}
)
