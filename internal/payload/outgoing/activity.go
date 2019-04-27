package outgoing

import "encoding/json"

type (
	// where the activity is and who is active
	Activity struct {
		UserID    uint64 `json:"userID,string"`
		Kind      string `json:"kind,omitempty"`
		MessageID uint64 `json:"messageID,string,omitempty"`
		ChannelID uint64 `json:"channelID,string,omitempty"`
		Present   bool   `json:"present"`
	}
)

func (p *Activity) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Activity: p})
}
