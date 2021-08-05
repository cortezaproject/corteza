package settings

type (
	Settings struct {
		LocalEnabled              bool
		SignupEnabled             bool
		EmailConfirmationRequired bool
		PasswordResetEnabled      bool
		ExternalEnabled           bool
		SplitCredentialsCheck     bool
		Providers                 []Provider
		MultiFactor               MultiFactor
	}

	MultiFactor struct {
		EmailOTP EmailOTP
		TOTP     TOTP
	}

	EmailOTP struct {
		// Can users use email for MFA
		Enabled bool

		// Is MFA with email enforced?
		Enforced bool
	}

	TOTP struct {
		// Can users use TOTP MFA?
		Enabled bool

		// Is TOTP MFA enforced?
		Enforced bool

		// TOTP issuer
		Issuer string
	}

	Provider struct {
		Handle      string
		Label       string
		IssuerUrl   string
		Key         string
		RedirectUrl string
		Secret      string
	}
)
