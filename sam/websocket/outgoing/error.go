package outgoing

type (
	Error struct {
		Message string `json:"m"`
	}
)

func (*Error) valid() bool { return true }
