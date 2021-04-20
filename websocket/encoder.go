package websocket

type (
	MessageEncoder interface {
		EncodeMessage() ([]byte, error)
	}
)
