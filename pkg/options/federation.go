package options

type (
	FederationOpt struct {
		Enabled bool   `env:"FEDERATION_ENABLED"`
		Host    string `env:"FEDERATION_HOST"`
	}
)

func Federation() (o *FederationOpt) {

	o = &FederationOpt{
		Enabled: false,
		Host:    `example.tld`,
	}

	fill(o)

	return
}
