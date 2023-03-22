package types

type (
	EmailNotification struct {
		To                []string
		Cc                []string
		ReplyTo           string
		Subject           string
		ContentPlain      string
		ContentHTML       string
		RemoteAttachments []string
	}
)
