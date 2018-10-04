package incoming

type (
	ExecCommand struct {
		ChannelID string            `json:"channelId"`
		Command   string            `json:"command"`
		Params    map[string]string `json:"params"`
		Input     string            `json:"input"`
	}
)
