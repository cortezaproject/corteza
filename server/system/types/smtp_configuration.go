package types

type (
	SmtpConfiguration struct {
		Host          string
		Port          uint
		Recipients    []string
		Username      string
		Password      string
		TLSInsecure   bool
		TLSServerName string
	}

	// SmtpCheckResult represents the messages returned after SMTP Host validation,
	// SMTP Server configurations check and Send test email process
	SmtpCheckResult struct {
		Host   string `json:"host"`
		Server string `json:"server"`
		Send   string `json:"send"`
	}
)
