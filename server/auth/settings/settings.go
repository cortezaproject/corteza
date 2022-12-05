package settings

type (
	Settings struct {
		LocalEnabled              bool
		SignupEnabled             bool
		EmailConfirmationRequired bool
		PasswordResetEnabled      bool
		PasswordCreateEnabled     bool
		ExternalEnabled           bool
		SplitCredentialsCheck     bool
		Providers                 []Provider
		Saml                      SAML
		MultiFactor               MultiFactor
	}

	SAML struct {
		Enabled bool

		// IdP name used on a login form
		Name string

		// SAML certificate
		Cert string

		// SAML certificate private key
		Key string

		// Sign AuthNRequest and assertion
		SignRequests bool

		// Signature method for signing
		SignMethod string

		// Post or redirect binding
		Binding string

		// Identity provider hostname
		IDP struct {
			URL string

			// identifier payload from idp
			IdentName       string
			IdentHandle     string
			IdentIdentifier string
		}
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
		Scope       string
		Usage       []string
	}
)

func (p Provider) HasUsage(u string) bool {
	for _, pu := range p.Usage {
		if pu == u {
			return true
		}
	}

	return false
}
