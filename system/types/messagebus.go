package types

type (
	QueueMessage struct {
		Queue   string `json:"queue"`
		Payload string `json:"payload"`
	}
)
