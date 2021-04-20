package outgoing

type (
	Payload struct {
		*Error `json:"error,omitempty"`

		*Prompt  `json:"prompt,omitempty"`
		*Prompts `json:"prompts,omitempty"`
	}

	// MessageEncoder This is same-same but different as using the json.Marshaler
	// (this one does not cause json.Marshal to call itself)
	MessageEncoder interface {
		EncodeMessage() ([]byte, error)
	}
)
