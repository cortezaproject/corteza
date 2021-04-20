package incoming

type Payload struct {
	// Token is JWT token provided by client as first message,
	// and will be passed whenever it changes
	*Token `json:"token"`
}
