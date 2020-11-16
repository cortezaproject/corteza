package options

type (
	FederationOpt struct {
		Enabled bool   `env:"FEDERATION_ENABLED"`
		Label   string `env:"FEDERATION_LABEL"`
		Host    string `env:"FEDERATION_HOST"`
	}
)

func Federation() (o *FederationOpt) {

	o = &FederationOpt{
		Enabled: false,
		Label:   `Example host`,
		Host:    `example.tld`,
	}

	fill(o)

	return
}
