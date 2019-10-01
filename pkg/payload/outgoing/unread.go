package outgoing

import "encoding/json"

type (
	Unread struct {
		ChannelID uint64 `json:"channelID,string,omitempty"`
		ThreadID  uint64 `json:"threadID,string,omitempty"`

		LastMessageID uint64 `json:"lastMessageID,string,omitempty"`
		Count         uint32 `json:"count"`

		ThreadCount uint32 `json:"threadCount"`
		ThreadTotal uint32 `json:"threadTotal,omitempty"`
	}
)

func (p *Unread) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Unread: p})
}
