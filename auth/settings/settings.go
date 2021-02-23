package settings

type (
	Settings struct {
		LocalEnabled              bool
		SignupEnabled             bool
		EmailConfirmationRequired bool
		PasswordResetEnabled      bool
		ExternalEnabled           bool
		Providers                 []Provider
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
